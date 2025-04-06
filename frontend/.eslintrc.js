const path = require("path");

module.exports = {
	parser: "@typescript-eslint/parser",
	parserOptions: {
		project: path.resolve(__dirname, "tsconfig.json"),
	},
	plugins: ["import", "@typescript-eslint"],
	extends: [
		"airbnb-typescript/base"
	],
	rules: {
		"@typescript-eslint/indent": ["error", "tab"],
		"@typescript-eslint/no-use-before-define": ["error", { functions: false, classes: false, variables: false }],
		"@typescript-eslint/quotes": ["error", "double"],
	},
};
