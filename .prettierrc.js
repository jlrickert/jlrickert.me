/** @type {import("prettier").Config} */
module.exports = {
	printWidth: 80,
	endOfLine: 'auto',
	singleQuote: true,
	tabWidth: 4,
	useTabs: true,
	overrides: [
		{
			files: ['*.yaml', '*.yml'],
			tabWidth: 2,
		},
		{
			files: '*.md',
			tabWidth: 2,
			proseWrap: true,
			useTabs: false,
		},
		{
			files: '*.svelte',
			options: { parser: 'svelte' },
		},
	],
};
