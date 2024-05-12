use std::{collections::HashMap, fs, io, path::Path};

fn main() -> Result<(), String> {
    let args = std::env::args().skip(1).collect::<Vec<_>>();
    if args.len() < 2 {
        panic!("Expects 2 arguments")
    }
    let file = args.get(0).unwrap();
    let output = args.get(1).unwrap();
    let showfile = args.get(2).map(|s| &**s).unwrap_or_default();
    let temp_dir = Path::new(file).with_extension("tmp");
    let data = fs::read(file).map_err(e("Reading MSI failed"))?;
    fs::create_dir(temp_dir.clone()).map_err(e("Temp dir already exists?"))?;
    if data[0..8] != [0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1] {
        println!("{file}");
        panic!("Invalid header!")
    }
    let version = u16::from_le_bytes([data[26], data[27]]);
    let sector_size = if version == 3 {
        if data[30..32] != [0x09, 0x00] {
            println!("{file}");
            panic!("Invalid shift")
        }
        512
    } else if version == 4 {
        if data[30..32] != [0x0c, 0x00] {
            println!("{file}");
            panic!("Invalid shift")
        }
        4096
    } else {
        println!("{file}");
        panic!("Invalid version")
    };
    let _num_dirs = get_u32(&data, 40);
    let _num_fats = get_u32(&data, 44);
    let mut dir_sector = get_u32(&data, 48);
    let mut minifat_sector = get_u32(&data, 60);
    let _num_mini = get_u32(&data, 64);
    let mut difat_entry = get_u32(&data, 68);
    let _num_difat = get_u32(&data, 72);
    let mut difat = vec![];
    for i in 0..109 {
        let difat_entry = get_u32(&data, 76 + 4 * i);
        if difat_entry != 0xFFFFFFFF {
            difat.push(difat_entry)
        }
    }
    while difat_entry != 0xFFFFFFFE {
        for i in 0..sector_size / 4 - 1 {
            let difat_entry_int = get_u32(&data, sector_size * (difat_entry as usize + 1) + i * 4);
            if difat_entry_int != 0xFFFFFFFF {
                difat.push(difat_entry_int)
            }
        }
        difat_entry = get_u32(&data, sector_size * (difat_entry as usize + 1));
    }
    let mut fat = vec![];
    for difat_sector in difat {
        for i in 0..sector_size / 4 {
            fat.push(get_u32(
                &data,
                sector_size * (difat_sector as usize + 1) + i * 4,
            ))
        }
    }
    let mut minifat = vec![];
    while minifat_sector != 0xFFFFFFFE {
        for i in 0..sector_size / 4 {
            minifat.push(get_u32(
                &data,
                sector_size * (minifat_sector as usize + 1) + i * 4,
            ));
        }
        minifat_sector = fat[minifat_sector as usize];
    }
    let mut mini_stream = vec![];
    let mut cabinets = vec![];
    while dir_sector != 0xFFFFFFFE {
        for dir_entry in 0..sector_size / 128 {
            let base_offset = sector_size * (dir_sector as usize + 1) + dir_entry * 128;
            let mut sector = get_u32(&data, base_offset + 116);
            if data[base_offset + 66] == 0x05 {
                while sector != 0xFFFFFFFE {
                    mini_stream.append(
                        &mut data[sector_size * (sector as usize + 1)
                            ..sector_size * (sector as usize + 2)]
                            .to_owned(),
                    );
                    sector = fat[sector as usize];
                }
                continue;
            }
            if data[base_offset + 66] == 0x01 || data[base_offset + 66] == 0x0 {
                continue;
            }
            let name_len = get_u16(&data, base_offset + 64);
            let mut name_wide = vec![];
            for i in 0..name_len / 2 {
                name_wide.push(get_u16(&data, base_offset + 2 * i as usize));
            }
            let entry_name = if name_wide.iter().all(|&c| c < 0x3800) {
                String::from_utf16(&name_wide)
                    .unwrap()
                    .trim_end_matches('\0')
                    .trim_start_matches('\x05')
                    .to_owned()
            } else {
                String::from_utf8(
                    name_wide
                        .iter()
                        .flat_map(|&chr| {
                            const MAP: &[u8] =
                                b"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz._";
                            if chr == 0 {
                                return vec![];
                            }
                            if chr == 0x4840 {
                                return vec![];
                            }
                            let chr = chr - 0x3800;
                            let chr1 = chr & 0x3F;
                            let chr2 = chr >> 6;
                            let mut bytes = vec![MAP[chr1 as usize]];
                            if chr2 == 0x40 {
                                return bytes;
                            }
                            bytes.push(MAP[chr2 as usize]);
                            bytes
                        })
                        .collect::<Vec<_>>(),
                )
                .unwrap()
            };
            let mut entry_size = get_u64(&data, base_offset + 120);
            if version == 3 {
                entry_size &= 0xffffffff;
            }
            let mut entry_data = vec![];
            if entry_size < 4096 {
                while sector != 0xFFFFFFFE {
                    entry_data.append(
                        &mut mini_stream[64 * sector as usize..64 * (sector as usize + 1)]
                            .to_owned(),
                    );
                    sector = minifat[sector as usize];
                }
            } else {
                while sector != 0xFFFFFFFE {
                    entry_data.append(
                        &mut data[sector_size * (sector as usize + 1)
                            ..sector_size * (sector as usize + 2)]
                            .to_owned(),
                    );
                    sector = fat[sector as usize];
                }
            }
            fs::write(
                temp_dir.join(&entry_name),
                &entry_data[0..entry_size as usize],
            )
            .map_err(e("Failed to write temp file"))?;
            if entry_name.contains(".cab") {
                let parent = std::path::Path::new(&temp_dir).parent().unwrap();
                fs::copy(temp_dir.join(&entry_name), parent.join(&entry_name))
                    .map_err(e("Failed to move cabinet"))?;
                cabinets.push(parent.join(&entry_name).to_string_lossy().into_owned());
            }
        }
        dir_sector = fat[dir_sector as usize];
    }
    // Now, extract & rename those files!
    let string_data =
        fs::read(temp_dir.join("_StringData")).map_err(e("Failed to read string data"))?;
    let string_pool =
        fs::read(temp_dir.join("_StringPool")).map_err(e("Failed to read string pool"))?;
    let mut offset = 0;
    let mut strings = vec![String::new()];
    for index in 1..string_pool.len() / 4 {
        let len = get_u16(&string_pool, index * 4);
        let data =
            String::from_utf8(string_data[offset..offset + len as usize].to_owned()).unwrap();
        strings.push(data);
        offset += len as usize;
    }
    let columns = fs::read(temp_dir.join("_Columns")).map_err(e("Failed to read column data"))?;
    let row_count = columns.len() / 8;
    let mut column_list = vec![];
    for i in 0..row_count {
        let table = get_u16(&columns, i * 2);
        let number = get_u16(&columns, i * 2 + row_count * 2);
        let name = get_u16(&columns, i * 2 + row_count * 4);
        let ty = get_u16(&columns, i * 2 + row_count * 6);
        column_list.push((table, number, name, ty));
    }
    let directory_table =
        fs::read(temp_dir.join("Directory")).map_err(e("Failed to read directory list"))?;
    let columns = get_columns("Directory", &strings, &column_list);
    let row_length = columns.last().unwrap().1 + columns.last().unwrap().2;
    let row_count = directory_table.len() / row_length as usize;
    let mut directories = HashMap::new();
    for i in 0..row_count {
        let directory_offs = get_val(&directory_table, i, 0, row_count, 2);
        let parent_offs = get_val(&directory_table, i, columns[1].1 as usize, row_count, 2);
        let dir_default_offs = get_val(&directory_table, i, columns[2].1 as usize, row_count, 2);
        directories.insert(
            strings[directory_offs as usize].as_str(),
            (
                strings[parent_offs as usize].as_str(),
                strings[dir_default_offs as usize].as_str(),
            ),
        );
    }
    let mut directory_map = HashMap::new();
    for (dir, (parent, name)) in directories.clone() {
        let mut full_path = String::new();
        let mut parent = parent;
        let mut name = name;
        while parent != "TARGET_DIR" && !parent.is_empty() {
            full_path.insert(0, '/');
            full_path.insert_str(0, name.split_once('|').unwrap_or((name, name)).1);
            (parent, name) = *directories.get(parent).unwrap();
        }
        directory_map.insert(dir, full_path);
    }
    let component_table =
        fs::read(temp_dir.join("Component")).map_err(e("Failed to read component list"))?;
    let columns = get_columns("Component", &strings, &column_list);
    let row_length = columns.last().unwrap().1 + columns.last().unwrap().2;
    let row_count = component_table.len() / row_length as usize;
    let mut components = HashMap::new();
    for i in 0..row_count {
        let component_offs = get_val(&component_table, i, 0, row_count, 2);
        let directory_offs = get_val(&component_table, i, columns[2].1 as usize, row_count, 2);
        let directory = strings[directory_offs as usize].as_str();
        components.insert(
            strings[component_offs as usize].as_str(),
            directory_map.get(directory).unwrap().as_str(),
        );
    }
    let file_table = fs::read(temp_dir.join("File")).map_err(e("Failed to read file list"))?;
    let columns = get_columns("File", &strings, &column_list);
    let row_length = columns.last().unwrap().1 + columns.last().unwrap().2;
    let row_count = file_table.len() / row_length as usize;
    let mut files = HashMap::new();
    for i in 0..row_count {
        let file_offs = get_val(&file_table, i, 0, row_count, 2);
        let filename_offs = get_val(&file_table, i, columns[2].1 as usize, row_count, 2);
        let filename = &strings[filename_offs as usize];
        let component_offs = get_val(&file_table, i, columns[1].1 as usize, row_count, 2);
        let component = &strings[component_offs as usize];
        let long_filename = filename
            .split_once('|')
            .map(|v| v.1)
            .unwrap_or(filename)
            .to_owned();
        files.insert(
            strings[file_offs as usize].as_str(),
            (*components.get(component.as_str()).unwrap()).to_owned() + long_filename.as_str(),
        );
    }
    if let Ok(media_table) = fs::read(temp_dir.join("Media")) {
        let columns = get_columns("Media", &strings, &column_list);
        let row_length = columns.last().unwrap().1 + columns.last().unwrap().2;
        let row_count = media_table.len() / row_length as usize;
        for i in 0..row_count {
            let cabinet = get_val(&media_table, i, columns[3].1 as usize, row_count, 2);
            let cab_name = strings[cabinet as usize].clone();
            if !cab_name.starts_with('#') && !cab_name.is_empty() {
                let path = Path::new(&temp_dir).parent().unwrap();
                cabinets.push(path.join(&cab_name).to_string_lossy().into_owned());
            }
        }
    }
    for cab in cabinets {
        println!("Extracting cabinet: {cab}");
        let base = Path::new(&cab).parent().unwrap();
        let cab_file = fs::File::open(&cab).map_err(e("Failed to open cabinet"))?;
        let mut cabinet = cab::Cabinet::new(cab_file).map_err(e("Failed to read cabinet"))?;
        for (path, name) in &files {
            if showfile == "-v" {
                println!("{name}");
            }
            if let Ok(mut reader) = cabinet.read_file(path) {
                fs::create_dir_all(base.join(output).join(name).parent().unwrap())
                    .map_err(e("Failed to create final directories"))?;
                let mut writer = fs::File::create(base.join(output).join(name))
                    .map_err(e("Failed to create file"))?;
                io::copy(&mut reader, &mut writer).map_err(e("Failed to copy file"))?;
            }
        }
    }
    fs::remove_dir_all(temp_dir).map_err(e("Failed to remove temp dir"))
}
fn get_u64(data: &[u8], offs: usize) -> u64 {
    u64::from_le_bytes([
        data[offs],
        data[offs + 1],
        data[offs + 2],
        data[offs + 3],
        data[offs + 4],
        data[offs + 5],
        data[offs + 6],
        data[offs + 7],
    ])
}
fn get_u32(data: &[u8], offs: usize) -> u32 {
    u32::from_le_bytes([data[offs], data[offs + 1], data[offs + 2], data[offs + 3]])
}

fn get_u16(data: &[u8], offs: usize) -> u16 {
    u16::from_le_bytes([data[offs], data[offs + 1]])
}

fn get_val(data: &[u8], row: usize, column_offs: usize, row_count: usize, size: usize) -> u32 {
    if size == 2 {
        get_u16(data, row * 2 + row_count * column_offs) as u32
    } else {
        get_u32(data, row * 2 + row_count * column_offs)
    }
}

fn get_columns(
    name: &str,
    strings: &[String],
    column_list: &[(u16, u16, u16, u16)],
) -> Vec<(u16, u16, u16)> {
    let mut table_columns = column_list
        .iter()
        .filter(|&v| strings[v.0 as usize] == name)
        .collect::<Vec<_>>();

    table_columns.sort_by_key(|&v| v.1);
    table_columns
        .iter()
        .scan(0, |i, &v| {
            let size = if v.3 & 0xefff == 0x0900 || v.3 & 0x0800 == 0x0800 {
                2
            } else {
                (v.3 & 0xff).max(2)
            };
            let old = *i;
            *i += size;
            Some((v.2, old, size))
        })
        .collect::<Vec<_>>()
}

fn e<T: std::fmt::Debug>(help: &str) -> impl Fn(T) -> String + '_ {
    move |e| format!("{help}: {e:?}")
}
