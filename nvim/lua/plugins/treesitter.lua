return {
	"nvim-treesitter/nvim-treesitter",
	lazy = false,
	build = ":TSUpdate",
	config = function()
		local config = require("nvim-treesitter")
		config.setup({
			ensure_installed = { "lua", "go", "javascript", "typescript", "yaml", "dockerfile" },
			highlight = { enable = true },
			indent = { enable = true },
		})
	end,
}
