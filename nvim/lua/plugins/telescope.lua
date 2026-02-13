return {
	{
		"nvim-telescope/telescope.nvim",
		version = "0.2.1",
		config = function()
			local builtin = require("telescope.builtin")
			vim.keymap.set(
				"n",
				"<C-p>",
				"<CMD>lua require'telescope.builtin'.find_files({ find_command = {'rg','--files','--hidden','-g','!.git',}})<CR>",
				{
					noremap = true,
				}
			)
			vim.keymap.set("n", "<leader>pw", function()
				builtin.grep_string({ search = vim.fn.input("Grep > ") })
			end)
			vim.keymap.set("n", "<leader>pp", builtin.live_grep, {})
			vim.keymap.set("n", "<C-g>", builtin.git_files, {})
		end,
	},
	{
		"nvim-telescope/telescope-ui-select.nvim",
		config = function()
			require("telescope").setup({
				extensions = {
					["ui-select"] = {
						require("telescope.themes").get_dropdown({
							-- even more opts
						}),
					},
				},
			})
			require("telescope").load_extension("ui-select")
		end,
	},
}
