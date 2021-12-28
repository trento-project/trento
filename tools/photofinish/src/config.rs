use std::fs;

#[derive(Debug)]
pub struct Scenario {
    pub label: String,
    pub files: Vec<String>,
}

pub fn get_config_file_content() -> String {
    match fs::read_to_string(".photofinish.toml".to_string()) {
        Ok(toml_content) => toml_content,
        Err(err) => {
            println!(
                "Error! Probably .photofinish.toml is missing\n{}",
                err.to_string()
            );
            String::new()
        }
    }
}

pub fn parse_scenarios(config: String) -> Vec<Scenario> {
    let toml_config: toml::value::Table = toml::from_str(&config).unwrap();
    toml_config
        .iter()
        .map(|(key, value)| {
            let scenario_files: Vec<String> = value["files"]
                .as_array()
                .unwrap()
                .iter()
                .map(|file_path| file_path.as_str().unwrap().to_string())
                .collect();
            Scenario {
                label: key.to_string(),
                files: scenario_files,
            }
        })
        .collect()
}
