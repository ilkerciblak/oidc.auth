return {
  {
    "nvim-neo-tree/neo-tree.nvim",
    branch = "v3.x",
    dependencies = {
      "nvim-lua/plenary.nvim",
      "MunifTanjim/nui.nvim",
      "nvim-tree/nvim-web-devicons", -- optional, but recommended
    },
    lazy = false,                 -- neo-tree will lazily load itself
    config = function()
      require("neo-tree").setup({
        close_if_last_window = true,
        window = {
          position = "right",
          width = 40,
        },
        filesystem = {
          filtered_items = {
            hide_by_name = {
              "node_modules",
              ".git",
            },
            hide_dotfiles = false,
            always_show_by_pattern = {
              ".env*"
            }
          },
        },
      })
      vim.keymap.set("n", "<C-n>", ":Neotree filesystem toggle position=right<CR>", {})
    end,
  },
}
