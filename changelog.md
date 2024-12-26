Added Which Command
Added a new "which" command to show the source location of prompts, helping users locate prompt files in their repositories.

- Added new which command that shows repository and file location for a given prompt
- Integrated which command into root command structure 

Added Edit Command
Added a new "edit" command to open prompts in the user's preferred editor ($EDITOR).

- Added new edit command that opens prompt files in the system editor
- Uses $EDITOR environment variable with fallback to vim
- Integrated edit command into root command structure 