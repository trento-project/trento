extern crate clap;

use clap::{App, Arg, SubCommand};

mod config;
mod list;
mod run;

#[tokio::main]
async fn main() {
    let options = App::new("photofinish")
        .version("1.0.0")
        .subcommand(SubCommand::with_name("list").about("list available event sets"))
        .subcommand(
            SubCommand::with_name("run")
                .about("injects a specific set of events")
                .arg(
                    Arg::with_name("url")
                        .short("u")
                        .long("url")
                        .default_value("http://localhost:8081/api/collect"),
                )
                .arg(
                    Arg::with_name("SET")
                        .help("name of the events set")
                        .required(true),
                ),
        )
        .get_matches();

    let http_client = reqwest::Client::new();

    let config = config::get_config_file_content();

    let scenarios = config::parse_scenarios(config);

    if let Some(_) = options.subcommand_matches("list") {
        list::show_list(&scenarios);
    }

    if let Some(run_options) = options.subcommand_matches("run") {
        let scenario_label = run_options.value_of("SET").unwrap();
        let endpoint_url = run_options.value_of("url").unwrap();
        run::run(
            endpoint_url,
            scenario_label.to_string(),
            scenarios,
            http_client,
        )
        .await;
    }
}
