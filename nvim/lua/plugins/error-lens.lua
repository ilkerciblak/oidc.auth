return {
	"chikko80/error-lens.nvim",
	config = function()
		require("error-lens").setup({
			enabled = true,
			prefix = 1,
		})
	end,
}
