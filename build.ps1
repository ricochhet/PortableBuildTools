Copy-Item -Path "./downloader.py" -Destination "./build/" -Verbose

cargo build --release
Copy-Item -Path "./target/release/rust-msiexec.exe" -Destination "./build/" -Verbose