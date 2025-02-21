import ora from "ora";

const spinner = ora({ text: "Cali Heys" }).start();

const run = () => {
	spinner.start();
};

run();
