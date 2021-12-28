use crate::config::Scenario;

pub fn show_list(scenarios: &Vec<Scenario>) -> () {
    for scenario in scenarios.iter() {
        print_scenario(scenario)
    }
}

fn print_scenario(scenario: &Scenario) -> () {
    println!(
        "NAME: {}\nFILES:\n{}\n",
        scenario.label,
        scenario.files.join("\n")
    )
}
