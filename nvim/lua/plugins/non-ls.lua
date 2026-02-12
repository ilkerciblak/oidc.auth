return {
	"nvimtools/none-ls.nvim",
	config = function()
		local ls = require("null-ls")
		ls.setup({
			sources = {
				ls.builtins.formatting.gofmt,
				ls.builtins.formatting.gofumpt,
				ls.builtins.formatting.goimports,
				ls.builtins.formatting.stylua,
				ls.builtins.diagnostics.golangci_lint,
				ls.builtins.completion.nvim_snippets,
				ls.builtins.hover.dictionary,
				ls.builtins.formatting.prettier,
			},
		})

		---@diagnostic disable-next-line: undefined-global
		vim.keymap.set("n", "<leader>ff", vim.lsp.buf.format, {})
	end,
}
