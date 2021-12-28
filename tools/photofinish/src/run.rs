use crate::config::Scenario;
use std::fs;

pub async fn run(
    remote_endpoint: &str,
    scenario_label: String,
    scenarios: Vec<Scenario>,
    http_client: reqwest::Client,
) -> () {
    let selected_scenario = scenarios
        .iter()
        .find(|current_scenario| current_scenario.label == scenario_label);

    match selected_scenario {
        None => println!("Non-existing scenario!"),
        Some(scenario) => {
            for file in scenario.files.iter() {
                let canonical_path = fs::canonicalize(file).unwrap();
                match fs::read_to_string(canonical_path) {
                    Ok(file_content) => {
                        let response = http_client
                            .post(remote_endpoint)
                            .body(file_content)
                            .send()
                            .await;
                        match response {
                            Ok(_) => {
                                println!("Successfully POSTed file: {}", file);
                            }
                            Err(_) => println!("Error while POSTing fixture:, {}", file),
                        }
                    }
                    Err(_) => println!("Couldn't read file: {}", file),
                }
            }
        }
    }
}
