-- This file  needs to have same structure as nvconfig.lua 
-- https://github.com/NvChad/NvChad/blob/v2.5/lua/nvconfig.lua

---@type ChadrcConfig
local M = {}

M.ui = {
	theme = "vscode_dark",

	nvdash = {
    load_on_startup = true,

    header = {
			"                                                                                ",
			"                                  ▒▓▓████████▓▓▒                                ",
			"                               ▒▓█████████████████▓                             ",
			"                             ▒▓█████████████████████▓                           ",
			"                            ▒████████▓▓▓▓▓▓▓██████████▒                         ",
			"                           ▒██████▓████████████▓███████                         ",
			"                           ███████▓▓▒▒▒▒▒▒▒▒▒▒▓▓███████▒                        ",
			"                          ▒████▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓████                        ",
			"                          ▓█▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓██                        ",
			"                          ▒▒▒▒▒▓▓▓▓▓▓▓▒▒▒▒▒▒▓▓▓▓▓▓▓▒▒▒▒▓                        ",
			"                          ▒██▓▓▒▒▓▓▓▓▓▓█████▓▓▓▓▓▒▒▓▓▓█▒▒                       ",
			"                         ▒▒▒█▒▒▒▒▒▓▓▒▒▒█▓▒█▓▒▒▓▓▒▒▒▒▒▒▓▒▒                       ",
			"                          ▓▒▓▒▒▒▒▒▒▒▒▒▓▓▒▒▒▓▒▒▒▒▒▒▒▒▒▓▒▓▒                       ",
			"                          ▓▒▒▓▓▓▓▓▓▓▓▓▒▒▒▒▒▒▓▓▓▓▓▓▓▓▓▓▒▓▒                       ",
			"                          ▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█▒                       ",
			"                          ▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█                        ",
			"                           █▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓                        ",
			"                           ▓▓▓▒▒▒▒▒▓▓▓▓▓███▓▓▓▓▒▒▒▒▒▒▓▓▒                        ",
			"                            █▓▓▒▒▒▓▓▓▒▓▒▒▒▒▓▓▓▓▓▒▒▒▓▓▓▓                         ",
			"                            ▒▓▓▓▓▓█▓▒▒▒▓▓▓▓▒▒▒▒█▓▓▓▓▓▓                          ",
			"                              ▓█▓███▒▒▒▒▓▓▓▒▒▒▓██▓▓▓▒                           ",
			"                               ▒▓███▓▓▓▓▓▓▓▓▓▓▓███▓                             ",
			"                                 ▒▓████████████▓▒                               ",
			"                                    ▒▓▓████▓▓▒                                  ",
			"                                                                                ",
      "     dd               iii tt                           VV     VV IIIII MM    MM ",
			"     dd rr rr    eee      tt      aa aa  gggggg   eee  VV     VV  III  MMM  MMM ",
      " dddddd rrr  r ee   e iii tttt   aa aaa gg   gg ee   e  VV   VV   III  MM MM MM ",
      "dd   dd rr     eeeee  iii tt    aa  aaa ggggggg eeeee    VV VV    III  MM    MM ",
	    " dddddd rr      eeeee iii  tttt  aaa aa      gg  eeeee    VVV    IIIII MM    MM ",
			"                                         ggggg                                  ",
 		},

    buttons = {
      { "  Find File", "SPC + ff", "Telescope find_files" },
      { "󰈚  Recent Files", "SPC + fo", "Telescope oldfiles" },
      { "󰈭  Find Word", "SPC + fw", "Telescope live_grep" },
      { "  Bookmarks", "SPC + ma", "Telescope marks" },
      { "  Themes", "SPC + th", "Telescope themes" },
      { "  Mappings", "SPC + ch", "NvCheatsheet" },
    },
  },

	-- hl_override = {
	-- 	Comment = { italic = true },
	-- 	["@comment"] = { italic = true },
	-- },
}

return M
