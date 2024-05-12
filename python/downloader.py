#!/usr/bin/env python3
# Originally forked from: https://gist.github.com/mmozeiko/7f3162ec2988e81e56d5c4e22cde9977
# Significant modifications have been made since to maximize compatibility.
# https://github.com/ricochhet/SdkStandalone

import io
import os
import sys
import stat
import json
import shutil
import hashlib
import zipfile
import logging
import tempfile
import argparse
import subprocess
import urllib.error
import urllib.request
from pathlib import Path

OUTPUT = Path("sdk_standalone")        # output folder
DOWNLOADS = Path("downloads") # temporary download files

HOST   = "x64" # or x86
DOWNLOAD_SPECTRE_LIBS = False # or False
DOWNLOAD_ARM_TARGETS = False # or False
MANIFEST_URL = "https://aka.ms/vs/17/release/channel"
MANIFEST_PREVIEW_URL = "https://aka.ms/vs/17/pre/channel"

# If there is a part of this downloader that you feel is missing by default, create an issue,
#   rather than switching download_all to true. . .
# https://github.com/ricochhet/SdkStandalone/issues
DOWNLOAD_ALL = False

# Only switch to 'True' if you know what you are doing.
OUTPUT_CLEANUP = False 

# DO NOT MODIFY
TARGETX64 = "x64"
TARGETX86 = "x86"
TARGETARM = "arm"
TARGETARM64 = "arm64"

ssl_context = None

def download(url):
  with urllib.request.urlopen(url, context=ssl_context) as res:
    return res.read()

total_download = 0

def download_progress(url, check, name, filename):
  fpath = DOWNLOADS / filename
  if fpath.exists():
    data = fpath.read_bytes()
    if hashlib.sha256(data).hexdigest() == check.lower():
      print(f"\r{name} ... OK")
      return data

  global total_download
  with fpath.open("wb") as f:
    data = io.BytesIO()
    with urllib.request.urlopen(url, context=ssl_context) as res:
      total = int(res.headers["Content-Length"])
      size = 0
      while True:
        block = res.read(1<<20)
        if not block:
          break
        f.write(block)
        data.write(block)
        size += len(block)
        perc = size * 100 // total
        print(f"\r{name} ... {perc}%", end="")
    print()
    data = data.getvalue()
    digest = hashlib.sha256(data).hexdigest()
    if check.lower() != digest:
      exit(f"Hash mismatch for f{pkg}")
    total_download += len(data)
    return data

# super crappy msi format parser just to find required .cab files
def get_msi_cabs(msi):
  index = 0
  while True:
    index = msi.find(b".cab", index+4)
    if index < 0:
      return
    yield msi[index-32:index+4].decode("ascii")

def first(items, cond):
  return next(item for item in items if cond(item))
  

### parse command-line arguments

ap = argparse.ArgumentParser()
ap.add_argument("--show-versions", action="store_true", help="Show available MSVC and Windows SDK versions")
ap.add_argument("--accept-license", action="store_true", help="Automatically accept license")
ap.add_argument("--msvc-version", help="Get specific MSVC version")
ap.add_argument("--sdk-version", help="Get specific Windows SDK version")
ap.add_argument("--preview", action="store_true", help="Use preview channel for Preview versions")

ap.add_argument("--dl-spectre", action="store_true", help="Download spectre mitigated libraries")
ap.add_argument("--dl-arm-targets", action="store_true", help="Download ARM/ARM64 targets")

args = ap.parse_args()

DOWNLOAD_SPECTRE_LIBS = True if args.dl_spectre else False
print(f"DOWNLOAD_SPECTRE_LIBS == {DOWNLOAD_SPECTRE_LIBS}")

DOWNLOAD_ARM_TARGETS = True if args.dl_arm_targets else False
print(f"DOWNLOAD_ARM_TARGETS == {DOWNLOAD_ARM_TARGETS}")

### get main manifest

URL = MANIFEST_PREVIEW_URL if args.preview else MANIFEST_URL

try:
  manifest = json.loads(download(URL))
except urllib.error.URLError as err:
  import ssl
  if isinstance(err.args[0], ssl.SSLCertVerificationError):
    # for more info about Python & issues with Windows certificates see https://stackoverflow.com/a/52074591
    print("ERROR: ssl certificate verification error")
    try:
      import certifi
    except ModuleNotFoundError:
      print("ERROR: please install 'certifi' package to use Mozilla certificates")
      print("ERROR: or update your Windows certs, see instructions here: https://woshub.com/updating-trusted-root-certificates-in-windows-10/#h2_3")
      exit()
    print("NOTE: retrying with certifi certificates")
    ssl_context = ssl.create_default_context(cafile=certifi.where())
    manifest = json.loads(download(URL))
  else:
    raise

### download VS manifest

ITEM_NAME = "Microsoft.VisualStudio.Manifests.VisualStudioPreview" if args.preview else "Microsoft.VisualStudio.Manifests.VisualStudio"

vs = first(manifest["channelItems"], lambda x: x["id"] == ITEM_NAME)
payload = vs["payloads"][0]["url"]

vsmanifest = json.loads(download(payload))


### find MSVC & WinSDK versions

packages = {}
for p in vsmanifest["packages"]:
  packages.setdefault(p["id"].lower(), []).append(p)

msvc = {}
sdk = {}

for pid,p in packages.items():
  if pid.startswith("Microsoft.VisualStudio.Component.VC.".lower()) and pid.endswith(".x86.x64".lower()):
    pver = ".".join(pid.split(".")[4:6])
    if pver[0].isnumeric():
      msvc[pver] = pid
  elif pid.startswith("Microsoft.VisualStudio.Component.Windows10SDK.".lower()) or \
      pid.startswith("Microsoft.VisualStudio.Component.Windows11SDK.".lower()):
    pver = pid.split(".")[-1]
    if pver.isnumeric():
      sdk[pver] = pid

if args.show_versions:
  print("MSVC versions:", " ".join(sorted(msvc.keys())))
  print("Windows SDK versions:", " ".join(sorted(sdk.keys())))
  exit(0)

msvc_ver = args.msvc_version or max(sorted(msvc.keys()))
sdk_ver = args.sdk_version or max(sorted(sdk.keys()))

if msvc_ver in msvc:
  msvc_pid = msvc[msvc_ver]
  msvc_ver = ".".join(msvc_pid.split(".")[4:-2])
else:
  exit(f"Unknown MSVC version: f{args.msvc_version}")

if sdk_ver in sdk:
  sdk_pid = sdk[sdk_ver]
else:
  exit(f"Unknown Windows SDK version: f{args.sdk_version}")

print(msvc_ver)
print(sdk_ver)
print(sdk_pid)

print(f"Downloading MSVC v{msvc_ver} and Windows SDK v{sdk_ver}")


### agree to license

tools = first(manifest["channelItems"], lambda x: x["id"] == "Microsoft.VisualStudio.Product.BuildTools")
resource = first(tools["localizedResources"], lambda x: x["language"] == "en-us")
license = resource["license"]

if not args.accept_license:
  accept = input(f"Do you accept Visual Studio license at {license} [Y/N] ? ")
  if not accept or accept[0].lower() != "y":
    exit(0)

OUTPUT.mkdir(exist_ok=True)
DOWNLOADS.mkdir(exist_ok=True)

logging.basicConfig(filename="downloader.log", level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

### download MSVC
if DOWNLOAD_ALL:
  prefixes = [f"microsoft.vc.{msvc_ver}", "microsoft.visualcpp", "microsoft.vs", "microsoft.visualstudio"]
  for prefix in prefixes:
      for pkg in packages:
          if pkg.startswith(prefix):
              try:
                p = first(packages[pkg], lambda p: "language" not in p or p.get("language") in (None, "en-US"))
              except StopIteration:
                p = None
              if p is not None and "payloads" in p:
                for payload in p["payloads"]:
                    filename = payload["fileName"]
                    download_progress(payload["url"], payload["sha256"], pkg, filename)
                    if zipfile.is_zipfile(DOWNLOADS / filename):
                      with zipfile.ZipFile(DOWNLOADS / filename) as z:
                          for name in z.namelist():
                              decoded_name = urllib.parse.unquote(name)
                              if decoded_name.startswith("Contents/"):
                                  out = OUTPUT / Path(decoded_name).relative_to("Contents")
                                  out.parent.mkdir(parents=True, exist_ok=True)
                                  out.write_bytes(z.read(name))
                                  logging.info(f"Package: {pkg}, Filename: {filename}, URL: {payload["url"]}, SHA256: {payload["sha256"]}, File: {str(out)}")
else:
  msvc_packages = [
    # MSVC vcvars
    "microsoft.visualstudio.vc.vcvars",
    "microsoft.visualstudio.vc.devcmd",
    "microsoft.visualcpp.tools.core.x86",
    f"microsoft.visualcpp.tools.host{HOST}.target{TARGETX64}",
    f"microsoft.visualcpp.tools.host{HOST}.target{TARGETX86}",
    # MSVC binaries x64
    f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETX64}.base",
    f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETX64}.res.base",
    # MSVC binaries x86
    f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETX86}.base",
    f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETX86}.res.base",
    # MSVC headers
    f"microsoft.vc.{msvc_ver}.crt.headers.base",
    # MSVC Libs x64
    f"microsoft.vc.{msvc_ver}.crt.{TARGETX64}.desktop.base",
    f"microsoft.vc.{msvc_ver}.crt.{TARGETX64}.store.base",
    # MSVC Libs x86
    f"microsoft.vc.{msvc_ver}.crt.{TARGETX86}.desktop.base",
    f"microsoft.vc.{msvc_ver}.crt.{TARGETX86}.store.base",
    # MSVC runtime source
    f"microsoft.vc.{msvc_ver}.crt.source.base",
    # ASAN
    f"microsoft.vc.{msvc_ver}.asan.headers.base",
    f"microsoft.vc.{msvc_ver}.asan.{TARGETX64}.base",
    f"microsoft.vc.{msvc_ver}.asan.{TARGETX86}.base",
    # MSVC tools
    f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETX64}.base",
    f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETX86}.base",
    # DIA SDK
    "microsoft.visualcpp.dia.sdk",
    # MSVC redist
    f"microsoft.visualcpp.crt.redist.{TARGETX64}",
    f"microsoft.visualcpp.crt.redist.{TARGETX86}",
    # MSVC UnitTest
    "microsoft.visualstudio.vc.unittest.desktop.build.core",
    "microsoft.visualstudio.testtools.codecoverage",
    # MSVC Cli
    f"microsoft.vc.{msvc_ver}.cli.{HOST}.base",
    # MSVC Store CRT
    f"microsoft.vc.{msvc_ver}.crt.{HOST}.store.base",
    # MSVC Common Tools
    "microsoft.visualcpp.tools.common.utils",
    # MSVC Log
    "microsoft.visualstudio.log",
    #f"microsoft.vc.{msvc_ver}.crt.redist.x64.base",
  ]

  if DOWNLOAD_SPECTRE_LIBS:
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.crt.{TARGETX64}.desktop.spectre.base")
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.crt.{TARGETX86}.desktop.spectre.base")

  if DOWNLOAD_ARM_TARGETS:
    msvc_packages.append(f"microsoft.visualcpp.tools.host{HOST}.target{TARGETARM}")
    msvc_packages.append(f"microsoft.visualcpp.tools.host{HOST}.target{TARGETARM64}")
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETARM}.base")
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETARM}.res.base")
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETARM64}.base")
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETARM64}.res.base")
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.crt.{TARGETARM}.desktop.base")
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.crt.{TARGETARM}.store.base")
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.crt.{TARGETARM64}.desktop.base")
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.crt.{TARGETARM64}.store.base")
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETARM}.base")
    msvc_packages.append(f"microsoft.vc.{msvc_ver}.tools.host{HOST}.target{TARGETARM64}.base")
    if DOWNLOAD_SPECTRE_LIBS:
      msvc_packages.append(f"microsoft.vc.{msvc_ver}.crt.{TARGETARM}.desktop.spectre.base")
      msvc_packages.append(f"microsoft.vc.{msvc_ver}.crt.{TARGETARM64}.desktop.spectre.base")    

  for pkg in msvc_packages:
    p = first(packages[pkg], lambda p: "language" not in p or p.get("language") in (None, "en-US"))
    for payload in p["payloads"]:
      filename = payload["fileName"]
      download_progress(payload["url"], payload["sha256"], pkg, filename)
      with zipfile.ZipFile(DOWNLOADS / filename) as z:
        for name in z.namelist():
          decoded_name = urllib.parse.unquote(name)
          if decoded_name.startswith("Contents/"):
            out = OUTPUT / Path(decoded_name).relative_to("Contents")
            out.parent.mkdir(parents=True, exist_ok=True)
            out.write_bytes(z.read(name))
            logging.info(f"Package: {pkg}, Filename: {filename}, URL: {payload["url"]}, SHA256: {payload["sha256"]}, File: {str(out)}")


### download Windows SDK

sdk_packages = [
  # Windows SDK tools (like rc.exe & mt.exe)
  f"Windows SDK for Windows Store Apps Tools-{TARGETX86}_en-us.msi",
  # Windows SDK headers
  f"Windows SDK for Windows Store Apps Headers-{TARGETX86}_en-us.msi",
  f"Windows SDK Desktop Headers {TARGETX86}-{TARGETX86}_en-us.msi",
  f"Windows SDK Desktop Headers {TARGETX64}-x86_en-us.msi",
  # Windows SDK libs
  f"Windows SDK for Windows Store Apps Libs-{TARGETX86}_en-us.msi",
  f"Windows SDK Desktop Libs {TARGETX64}-{TARGETX86}_en-us.msi",
  f"Windows SDK Desktop Libs {TARGETX86}-{TARGETX86}_en-us.msi",
  # Windows SDK tools
  f"Windows SDK Desktop Tools {TARGETX64}-{TARGETX86}_en-us.msi",
  f"Windows SDK Desktop Tools {TARGETX86}-{TARGETX86}_en-us.msi",
  # CRT headers & libs
  f"Universal CRT Headers Libraries and Sources-{TARGETX86}_en-us.msi",
  # CRT redist
  f"Universal CRT Redistributable-{TARGETX86}_en-us.msi",
  # Signing tools
  f"Windows SDK Signing Tools-{TARGETX86}_en-us.msi",
]

if DOWNLOAD_ARM_TARGETS:
  sdk_packages.append(f"Windows SDK Desktop Headers {TARGETARM64}-{TARGETX86}_en-us.msi")
  sdk_packages.append(f"Windows SDK Desktop Headers {TARGETARM}-{TARGETX86}_en-us.msi")
  sdk_packages.append(f"Windows SDK Desktop Libs {TARGETARM64}-{TARGETX86}_en-us.msi")
  sdk_packages.append(f"Windows SDK Desktop Libs {TARGETARM}-{TARGETX86}_en-us.msi")
  sdk_packages.append(f"Windows SDK ARM Desktop Tools-{TARGETX86}_en-us.msi")
  sdk_packages.append(f"Windows SDK Desktop Tools {TARGETARM64}-{TARGETX86}_en-us.msi")

with tempfile.TemporaryDirectory(dir=DOWNLOADS) as d:
  dstX64 = Path(d)

  sdk_pkg = packages[sdk_pid][0]
  sdk_pkg = packages[first(sdk_pkg["dependencies"], lambda x: True).lower()][0]

  msi = []
  cabs = []

  # download msi files
  for pkg in sdk_packages:
    payload = first(sdk_pkg["payloads"], lambda p: p["fileName"] == f"Installers\\{pkg}") #NOSONAR
    msi.append(DOWNLOADS / pkg)
    data = download_progress(payload["url"], payload["sha256"], pkg, pkg)
    cabs += list(get_msi_cabs(data))
    logging.info(f"Package: {pkg}, Filename: {payload}, URL: {payload["url"]}, SHA256: {payload["sha256"]}, File: {str(pkg)}")

  # download .cab files
  for pkg in cabs:
    payload = first(sdk_pkg["payloads"], lambda p: p["fileName"] == f"Installers\\{pkg}") #NOSONAR
    download_progress(payload["url"], payload["sha256"], pkg, pkg)
    logging.info(f"Package: {pkg}, Filename: {payload}, URL: {payload["url"]}, SHA256: {payload["sha256"]}, File: {str(pkg)}")

  print("Unpacking msi files...")

  # run msi installers
  for m in msi:
    subprocess.check_call(["MSIExtract", f"{m}", f"{OUTPUT.resolve()}"])


### versions

msvcv = list((OUTPUT / "VC/Tools/MSVC").glob("*"))[0].name #NOSONAR
sdkv = list((OUTPUT / "Windows Kits/10/bin").glob("*"))[0].name


# place debug CRT runtime files into MSVC folder (not what real Visual Studio installer does... but is reasonable)

dstX64 = OUTPUT / "VC/Tools/MSVC" / msvcv / f"bin/Host{HOST}/{TARGETX64}"
dstX86 = OUTPUT / "VC/Tools/MSVC" / msvcv / f"bin/Host{HOST}/{TARGETX86}"
if DOWNLOAD_ARM_TARGETS:
  dstARM = OUTPUT / "VC/Tools/MSVC" / msvcv / f"bin/Host{HOST}/{TARGETARM}"
  dstARM64 = OUTPUT / "VC/Tools/MSVC" / msvcv / f"bin/Host{HOST}/{TARGETARM64}"

DOWNLOAD_FOLDER = Path("crtd")
(DOWNLOADS / DOWNLOAD_FOLDER).mkdir(exist_ok=True)

pkg = "microsoft.visualcpp.runtimedebug.14"
dbg = first(packages[pkg], lambda p: p["chip"] == HOST)
for payload in dbg["payloads"]:
  name = payload["fileName"]
  download_progress(payload["url"], payload["sha256"], name, DOWNLOAD_FOLDER / name)

msi = DOWNLOADS / DOWNLOAD_FOLDER / first(dbg["payloads"], lambda p: p["fileName"].endswith(".msi"))["fileName"]

with tempfile.TemporaryDirectory(dir=OUTPUT) as d2:
  subprocess.check_call(["MSIExtract", f"{str(msi)}", f"{Path(d2).resolve()}"])
  for f in first(Path(d2).glob("System*"), lambda x: True).iterdir():
    shutil.copy(f, dstX64 / f.name)
    shutil.copy(f, dstX86 / f.name)
    if DOWNLOAD_ARM_TARGETS:
      shutil.copy(f, dstARM / f.name)
      shutil.copy(f, dstARM64 / f.name)


# download DIA SDK and put msdia140.dll file into MSVC folder

DOWNLOAD_FOLDER = Path("dia")
(DOWNLOADS / DOWNLOAD_FOLDER).mkdir(exist_ok=True)

pkg = "microsoft.visualc.140.dia.sdk.msi"
dia = packages[pkg][0]
for payload in dia["payloads"]:
  name = payload["fileName"]
  download_progress(payload["url"], payload["sha256"], name, DOWNLOAD_FOLDER / name)

msi = DOWNLOADS / DOWNLOAD_FOLDER / first(dia["payloads"], lambda p: p["fileName"].endswith(".msi"))["fileName"]

with tempfile.TemporaryDirectory(dir=DOWNLOADS) as d2:
  subprocess.check_call(["MSIExtract", f"{str(msi)}", f"{Path(d2).resolve()}"])

  if HOST == TARGETX86: msdia = "msdia140.dll" #NOSONAR
  elif HOST == TARGETX64: msdia = "amd64/msdia140.dll"
  else: exit("unknown")

  # remove read-only attribute
  targetX64 = dstX64 / "msdia140.dll"
  if targetX64.exists():
    targetX64.chmod(stat.S_IWRITE)

  targetX86 = dstX86 / "msdia140.dll"
  if targetX86.exists():
    targetX86.chmod(stat.S_IWRITE)

  if DOWNLOAD_ARM_TARGETS:
    targetARM = dstARM / "msdia140.dll"
    if targetARM.exists():
      targetARM.chmod(stat.S_IWRITE)

    targetARM64 = dstARM64 / "msdia140.dll"
    if targetARM64.exists():
      targetARM64.chmod(stat.S_IWRITE)

  src = Path(d2) / "Program Files/Microsoft Visual Studio 14.0/DIA SDK/bin" / msdia
  shutil.copy(src, targetX64)
  shutil.copy(src, targetX86)
  if DOWNLOAD_ARM_TARGETS:
    shutil.copy(src, targetARM)
    shutil.copy(src, targetARM64)


### cleanup

if OUTPUT_CLEANUP:
  shutil.rmtree(OUTPUT / "Common7", ignore_errors=True)
  for f in ["Auxiliary", f"lib/{TARGETX64}/store", f"lib/{TARGETX64}/uwp"]:
    shutil.rmtree(OUTPUT / "VC/Tools/MSVC" / msvcv / f)
  for f in OUTPUT.glob("*.msi"):
    f.unlink()
  for f in ["Catalogs", "DesignTime", f"bin/{sdkv}/chpe", f"Lib/{sdkv}/ucrt_enclave"]:
    shutil.rmtree(OUTPUT / "Windows Kits/10" / f, ignore_errors=True)
  for arch in [TARGETX86, TARGETX64, TARGETARM, TARGETARM64]:
    if arch != TARGETX64:
      shutil.rmtree(OUTPUT / "Windows Kits/10/Lib" / sdkv / "ucrt" / arch)
      shutil.rmtree(OUTPUT / "Windows Kits/10/Lib" / sdkv / "um" / arch)
    if arch != HOST:
      shutil.rmtree(OUTPUT / "VC/Tools/MSVC" / msvcv / f"bin/Host{arch}", ignore_errors=True)
      shutil.rmtree(OUTPUT / "Windows Kits/10/bin" / sdkv / arch)

for arch in [TARGETX86, TARGETX64, TARGETARM, TARGETARM64]:
  if arch != HOST:
    shutil.rmtree(OUTPUT / "VC/Tools/MSVC" / msvcv / f"bin/Host{arch}", ignore_errors=True)

# executable that is collecting & sending telemetry every time cl/link runs
(OUTPUT / "VC/Tools/MSVC" / msvcv / f"bin/Host{HOST}/{TARGETX64}/vctip.exe").unlink(missing_ok=True)
(OUTPUT / "VC/Tools/MSVC" / msvcv / f"bin/Host{HOST}/{TARGETX86}/vctip.exe").unlink(missing_ok=True)
if DOWNLOAD_ARM_TARGETS:
  (OUTPUT / "VC/Tools/MSVC" / msvcv / f"bin/Host{HOST}/{TARGETARM}/vctip.exe").unlink(missing_ok=True)
  (OUTPUT / "VC/Tools/MSVC" / msvcv / f"bin/Host{HOST}/{TARGETARM64}/vctip.exe").unlink(missing_ok=True)

SET_VARS64 = f"""@echo off

SET ROOT=%~dp0
SET VSINSTALLDIR=%ROOT%\\
SET VCINSTALLDIR=%VSINSTALLDIR%VC\\
SET VS140COMNTOOLS=%VSINSTALLDIR%Common7\\Tools\\
SET UCRTVersion={sdkv}
SET WindowsSdkDir=%VSINSTALLDIR%Windows Kits\\10\\
SET UniversalCRTSdkDir=%WindowsSdkDir%
SET WindowsSDKVersion={sdkv}\\
SET WindowsSDKLibVersion={sdkv}\\
SET WindowsSDK_ExecutablePath_x64=%VSINSTALLDIR%Windows Kits\\10\\BIN\\%WindowsSDKVersion%{TARGETX64}\\
SET VSSPECTRELIBS=false

SET LIB=
SET INCLUDE=
SET LIBPATH=

@if exist "%VSINSTALLDIR%Common7\\Tools" set PATH=%VSINSTALLDIR%Common7\\Tools;%PATH%
@if exist "%VSINSTALLDIR%Common7\\IDE" set PATH=%VSINSTALLDIR%Common7\\IDE;%PATH%
@if exist "%VCINSTALLDIR%VCPackages" set PATH=%VCINSTALLDIR%VCPackages;%PATH%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETX64}" set PATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETX64};%PATH%

@if not "%UCRTVersion%" == "" @set INCLUDE=%UniversalCRTSdkDir%include\\%UCRTVersion%\\ucrt;%INCLUDE%
@if not "%UCRTVersion%" == "" @set LIB=%UniversalCRTSdkDir%lib\\%UCRTVersion%\\ucrt\\{TARGETX64};%LIB%

@if not "%WindowsSdkDir%" == "" @set PATH=%WindowsSdkDir%BIN\\{sdkv}\\{TARGETX64};%WindowsSdkDir%BIN\\{sdkv}\\{TARGETX86};%PATH%
@if not "%WindowsSdkDir%" == "" @set INCLUDE=%WindowsSdkDir%include\\%WindowsSDKVersion%shared;%WindowsSdkDir%include\\%WindowsSDKVersion%um;%WindowsSdkDir%include\\%WindowsSDKVersion%winrt;%INCLUDE%
@if not "%WindowsSdkDir%" == "" @set LIB=%WindowsSdkDir%lib\\%WindowsSDKLibVersion%um\\{TARGETX64};%LIB%
@if not "%WindowsSdkDir%" == "" @set LIBPATH=%WindowsLibPath%;%ExtensionSDKDir%\\Microsoft.VCLibs\\14.0\\References\\CommonConfiguration\\neutral;%LIBPATH%

@if not "%WindowsSDK_ExecutablePath_x64%" == "" @set PATH=%WindowsSDK_ExecutablePath_x64%;%PATH%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX64}\\store" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX64}\\store;%LIB%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\INCLUDE" set INCLUDE=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\INCLUDE;%INCLUDE%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\INCLUDE" set INCLUDE=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\INCLUDE;%INCLUDE%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETX64}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETX64};%LIB%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX64}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX64};%LIB%
if %VSSPECTRELIBS% == true (
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETX64}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETX64};%LIB%
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETX64}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETX64};%LIB%
)

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETX64}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETX64};%LIBPATH%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX64}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX64};%LIBPATH%
if %VSSPECTRELIBS% == true (
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETX64}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETX64};%LIBPATH%
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETX64}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETX64};%LIBPATH%
)

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX64}\\store" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX64}\\store;%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX64}\\store\\references;%LIBPATH%
"""

SET_VARS_ARM64 = f"""@echo off

SET ROOT=%~dp0
SET VSINSTALLDIR=%ROOT%\\
SET VCINSTALLDIR=%VSINSTALLDIR%VC\\
SET VS140COMNTOOLS=%VSINSTALLDIR%Common7\\Tools\\
SET UCRTVersion={sdkv}
SET WindowsSdkDir=%VSINSTALLDIR%Windows Kits\\10\\
SET UniversalCRTSdkDir=%WindowsSdkDir%
SET WindowsSDKVersion={sdkv}\\
SET WindowsSDKLibVersion={sdkv}\\
SET WindowsSDK_ExecutablePath_x64=%VSINSTALLDIR%Windows Kits\\10\\BIN\\%WindowsSDKVersion%{TARGETARM64}\\
SET VSSPECTRELIBS=false

SET LIB=
SET INCLUDE=
SET LIBPATH=

@if exist "%VSINSTALLDIR%Common7\\Tools" set PATH=%VSINSTALLDIR%Common7\\Tools;%PATH%
@if exist "%VSINSTALLDIR%Common7\\IDE" set PATH=%VSINSTALLDIR%Common7\\IDE;%PATH%
@if exist "%VCINSTALLDIR%VCPackages" set PATH=%VCINSTALLDIR%VCPackages;%PATH%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETARM64}" set PATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETARM64};%PATH%

@if not "%UCRTVersion%" == "" @set INCLUDE=%UniversalCRTSdkDir%include\\%UCRTVersion%\\ucrt;%INCLUDE%
@if not "%UCRTVersion%" == "" @set LIB=%UniversalCRTSdkDir%lib\\%UCRTVersion%\\ucrt\\{TARGETARM64};%LIB%

@if not "%WindowsSdkDir%" == "" @set PATH=%WindowsSdkDir%BIN\\{sdkv}\\{TARGETARM64};%WindowsSdkDir%BIN\\{sdkv}\\{TARGETARM64};%PATH%
@if not "%WindowsSdkDir%" == "" @set INCLUDE=%WindowsSdkDir%include\\%WindowsSDKVersion%shared;%WindowsSdkDir%include\\%WindowsSDKVersion%um;%WindowsSdkDir%include\\%WindowsSDKVersion%winrt;%INCLUDE%
@if not "%WindowsSdkDir%" == "" @set LIB=%WindowsSdkDir%lib\\%WindowsSDKLibVersion%um\\{TARGETARM64};%LIB%
@if not "%WindowsSdkDir%" == "" @set LIBPATH=%WindowsLibPath%;%ExtensionSDKDir%\\Microsoft.VCLibs\\14.0\\References\\CommonConfiguration\\neutral;%LIBPATH%

@if not "%WindowsSDK_ExecutablePath_x64%" == "" @set PATH=%WindowsSDK_ExecutablePath_x64%;%PATH%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM64}\\store" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM64}\\store;%LIB%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\INCLUDE" set INCLUDE=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\INCLUDE;%INCLUDE%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\INCLUDE" set INCLUDE=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\INCLUDE;%INCLUDE%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETARM64}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETARM64};%LIB%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM64}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM64};%LIB%
if %VSSPECTRELIBS% == true (
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETARM64}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETARM64};%LIB%
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETARM64}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETARM64};%LIB%
)

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETARM64}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETARM64};%LIBPATH%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM64}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM64};%LIBPATH%
if %VSSPECTRELIBS% == true (
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETARM64}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETARM64};%LIBPATH%
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETARM64}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETARM64};%LIBPATH%
)

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM64}\\store" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM64}\\store;%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM64}\\store\\references;%LIBPATH%
"""

SET_VARS32 = f"""@echo off

SET ROOT=%~dp0
SET VSINSTALLDIR=%ROOT%\\
SET VCINSTALLDIR=%VSINSTALLDIR%VC\\
SET VS140COMNTOOLS=%VSINSTALLDIR%Common7\\Tools\\
SET UCRTVersion={sdkv}
SET WindowsSdkDir=%VSINSTALLDIR%Windows Kits\\10\\
SET UniversalCRTSdkDir=%WindowsSdkDir%
SET WindowsSDKVersion={sdkv}\\
SET WindowsSDKLibVersion={sdkv}\\
SET WindowsSDK_ExecutablePath_x64=%VSINSTALLDIR%Windows Kits\\10\\BIN\\%WindowsSDKVersion%{TARGETX86}\\
SET VSSPECTRELIBS=false

SET LIB=
SET INCLUDE=
SET LIBPATH=

@if exist "%VSINSTALLDIR%Common7\\Tools" set PATH=%VSINSTALLDIR%Common7\\Tools;%PATH%
@if exist "%VSINSTALLDIR%Common7\\IDE" set PATH=%VSINSTALLDIR%Common7\\IDE;%PATH%
@if exist "%VCINSTALLDIR%VCPackages" set PATH=%VCINSTALLDIR%VCPackages;%PATH%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETX64}" set PATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETX64};%PATH%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETX86}" set PATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETX86};%PATH%

@if not "%UCRTVersion%" == "" @set INCLUDE=%UniversalCRTSdkDir%include\\%UCRTVersion%\\ucrt;%INCLUDE%
@if not "%UCRTVersion%" == "" @set LIB=%UniversalCRTSdkDir%lib\\%UCRTVersion%\\ucrt\\{TARGETX86};%LIB%

@if not "%WindowsSdkDir%" == "" @set PATH=%WindowsSdkDir%BIN\\{sdkv}\\{TARGETX86};%WindowsSdkDir%BIN\\{sdkv}\\{TARGETX86};%PATH%
@if not "%WindowsSdkDir%" == "" @set INCLUDE=%WindowsSdkDir%include\\%WindowsSDKVersion%shared;%WindowsSdkDir%include\\%WindowsSDKVersion%um;%WindowsSdkDir%include\\%WindowsSDKVersion%winrt;%INCLUDE%
@if not "%WindowsSdkDir%" == "" @set LIB=%WindowsSdkDir%lib\\%WindowsSDKLibVersion%um\\{TARGETX86};%LIB%
@if not "%WindowsSdkDir%" == "" @set LIBPATH=%WindowsLibPath%;%ExtensionSDKDir%\\Microsoft.VCLibs\\14.0\\References\\CommonConfiguration\\neutral;%LIBPATH%

@if not "%WindowsSDK_ExecutablePath_x64%" == "" @set PATH=%WindowsSDK_ExecutablePath_x64%;%PATH%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\INCLUDE" set INCLUDE=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\INCLUDE;%INCLUDE%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\INCLUDE" set INCLUDE=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\INCLUDE;%INCLUDE%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETX86}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETX86};%LIB%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX86}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX86};%LIB%
if %VSSPECTRELIBS% == true (
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETX86}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETX86};%LIB%
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETX86}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETX86};%LIB%
)

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX86}\\store" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX86}\\store;%LIB%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETX86}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETX86};%LIBPATH%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX86}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX86};%LIBPATH%
if %VSSPECTRELIBS% == true (
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETX86}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETX86};%LIBPATH%
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETX86}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETX86};%LIBPATH%
)

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX86}\\store" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX86}\\store;%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETX86}\\store\\references;%LIBPATH%
"""

SET_VARS_ARM32 = f"""@echo off

SET ROOT=%~dp0
SET VSINSTALLDIR=%ROOT%\\
SET VCINSTALLDIR=%VSINSTALLDIR%VC\\
SET VS140COMNTOOLS=%VSINSTALLDIR%Common7\\Tools\\
SET UCRTVersion={sdkv}
SET WindowsSdkDir=%VSINSTALLDIR%Windows Kits\\10\\
SET UniversalCRTSdkDir=%WindowsSdkDir%
SET WindowsSDKVersion={sdkv}\\
SET WindowsSDKLibVersion={sdkv}\\
SET WindowsSDK_ExecutablePath_x64=%VSINSTALLDIR%Windows Kits\\10\\BIN\\%WindowsSDKVersion%{TARGETARM}\\
SET VSSPECTRELIBS=false

SET LIB=
SET INCLUDE=
SET LIBPATH=

@if exist "%VSINSTALLDIR%Common7\\Tools" set PATH=%VSINSTALLDIR%Common7\\Tools;%PATH%
@if exist "%VSINSTALLDIR%Common7\\IDE" set PATH=%VSINSTALLDIR%Common7\\IDE;%PATH%
@if exist "%VCINSTALLDIR%VCPackages" set PATH=%VCINSTALLDIR%VCPackages;%PATH%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETARM64}" set PATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETARM64};%PATH%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETARM}" set PATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\BIN\\Host{HOST}\\{TARGETARM};%PATH%

@if not "%UCRTVersion%" == "" @set INCLUDE=%UniversalCRTSdkDir%include\\%UCRTVersion%\\ucrt;%INCLUDE%
@if not "%UCRTVersion%" == "" @set LIB=%UniversalCRTSdkDir%lib\\%UCRTVersion%\\ucrt\\{TARGETARM};%LIB%

@if not "%WindowsSdkDir%" == "" @set PATH=%WindowsSdkDir%BIN\\{sdkv}\\{TARGETARM};%WindowsSdkDir%BIN\\{sdkv}\\{TARGETARM};%PATH%
@if not "%WindowsSdkDir%" == "" @set INCLUDE=%WindowsSdkDir%include\\%WindowsSDKVersion%shared;%WindowsSdkDir%include\\%WindowsSDKVersion%um;%WindowsSdkDir%include\\%WindowsSDKVersion%winrt;%INCLUDE%
@if not "%WindowsSdkDir%" == "" @set LIB=%WindowsSdkDir%lib\\%WindowsSDKLibVersion%um\\{TARGETARM};%LIB%
@if not "%WindowsSdkDir%" == "" @set LIBPATH=%WindowsLibPath%;%ExtensionSDKDir%\\Microsoft.VCLibs\\14.0\\References\\CommonConfiguration\\neutral;%LIBPATH%

@if not "%WindowsSDK_ExecutablePath_x64%" == "" @set PATH=%WindowsSDK_ExecutablePath_x64%;%PATH%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\INCLUDE" set INCLUDE=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\INCLUDE;%INCLUDE%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\INCLUDE" set INCLUDE=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\INCLUDE;%INCLUDE%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETARM}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETARM};%LIB%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM};%LIB%
if %VSSPECTRELIBS% == true (
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETARM}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETARM};%LIB%
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETARM}" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETARM};%LIB%
)

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM}\\store" set LIB=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM}\\store;%LIB%

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETARM}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\{TARGETARM};%LIBPATH%
@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM};%LIBPATH%
if %VSSPECTRELIBS% == true (
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETARM}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\ATLMFC\\LIB\\SPECTRE\\{TARGETARM};%LIBPATH%
  @if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETARM}" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\SPECTRE\\{TARGETARM};%LIBPATH%
)

@if exist "%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM}\\store" set LIBPATH=%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM}\\store;%VCINSTALLDIR%\\Tools\\MSVC\\{msvcv}\\LIB\\{TARGETARM}\\store\\references;%LIBPATH%
"""

(OUTPUT / "set_vars32.bat").write_text(SET_VARS32)
(OUTPUT / "set_vars64.bat").write_text(SET_VARS64)
(OUTPUT / "set_vars_arm32.bat").write_text(SET_VARS_ARM32)
(OUTPUT / "set_vars_arm64.bat").write_text(SET_VARS_ARM64)

print(f"Total downloaded: {total_download>>20} MB")
print("Done!")
