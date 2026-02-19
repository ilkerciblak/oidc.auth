return {
	{
			"catppuccin/nvim",
			name = "catppuccin",
			priority = 1000,
			config = function()
				require("catppuccin").setup({

        lsp_styles ={
          inlay_hints = {
            background= true,
          },
        },
        auto_intergrations = true,
        })
				vim.cmd.colorscheme("catppuccin-mocha")
			end,
		},
}
