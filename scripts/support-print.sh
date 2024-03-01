#!/usr/bin/env bash

source $(dirname $0)/support-mecho.sh


THEME_PRIMARY_COLOR="yellow"
THEME_SECONDARY_COLOR="cyan"
THEME_NUMBER_OF_COLUMNS="80"
THEME_TITLE_LINE_SYMBOL="="
THEME_SECTIOM_LINE_SYMBOL="-"

# Prints a title as shown below
#
# ================================================================================
#
#                               INSTALL DEPENDENCIES
#
# ================================================================================
print::title() {
    local PADDING=$(((${THEME_NUMBER_OF_COLUMNS}/2)-(`echo "$1" | wc -c`)/2))
    local SYMBOL=";"

    local MESSAGE=$(printf "%${PADDING}s%s%-${PADDING}s\n" "$SYMBOL" "$1" "$SYMBOL" | tr ' ;' ' ' |  tr "[:lower:]" "[:upper:]")
    local LINE=$(printf "%${THEME_NUMBER_OF_COLUMNS}s" "$SYMBOL"| tr ' ;' "${THEME_TITLE_LINE_SYMBOL}")

    mecho "\n[${THEME_PRIMARY_COLOR}][underline]${LINE}[/]\n"
    mecho "[${THEME_PRIMARY_COLOR}][bold]${MESSAGE}[/]\n"
    mecho "[${THEME_PRIMARY_COLOR}][underline]${LINE}[/]\n"
}

# Prints a section as shown below:
# 
# --------------------------------------------------------------------------------
#
#                               Installing Homebrew
#
# --------------------------------------------------------------------------------
print::section() {
    local PADDING=$(((${THEME_NUMBER_OF_COLUMNS}/2)-(`echo "$1" | wc -c`)/2))
    local SYMBOL=";"

    local MESSAGE=$(printf "%${PADDING}s%s%-${PADDING}s\n" "$SYMBOL" "$1" "$SYMBOL" | tr ' ;' ' ')
    local LINE=$(printf "%${THEME_NUMBER_OF_COLUMNS}s" "$SYMBOL"| tr ' ;' "${THEME_SECTIOM_LINE_SYMBOL}")

    mecho "\n[${THEME_SECONDARY_COLOR}][underline]${LINE}[/]\n"
    mecho "[${THEME_SECONDARY_COLOR}][bold]${MESSAGE}[/]\n"
    mecho "[${THEME_SECONDARY_COLOR}][underline]${LINE}[/]\n"
}

# Prints a title paragraph respecting the number of columns and breaking
# lines if needed.
print::title_paragraph() {
    local MESSAGE=$(echo "$1" | fold -w ${THEME_NUMBER_OF_COLUMNS} -s)
    mecho "[${THEME_PRIMARY_COLOR}][italic]${MESSAGE}[/]\n"
}

# Prints a section paragraph respecting the number of columns and breaking
# lines if needed.
print::section_paragraph() {
    local MESSAGE=$(echo "$1" | fold -w ${THEME_NUMBER_OF_COLUMNS} -s)
    mecho "[${THEME_SECONDARY_COLOR}][italic]${MESSAGE}[/]\n"
}

# Prints an info message respecting the number of columns and breaking
# lines if needed.
print::info() {
    local MESSAGE=$(echo "$1" | fold -w ${THEME_NUMBER_OF_COLUMNS} -s)
    mecho "\n[yellow][bold][info][/] [yellow][italic]${MESSAGE}[/]"
}

# Prints a success message respecting the number of columns and breaking
# lines if needed.
print::success() {
    local MESSAGE=$(echo "$1" | fold -w ${THEME_NUMBER_OF_COLUMNS} -s)
    mecho "\n[green][bold][ok][/] [green][italic]${MESSAGE}[/]"
}

# Prints an error respecting the number of columns and breaking
# lines if needed.
print::error() {
    local MESSAGE=$(echo "$1" | fold -w ${THEME_NUMBER_OF_COLUMNS} -s)
    mecho "\n[red][bold][error][/] [red][italic]${MESSAGE}[/]"
}

# Prints & executes a command message respecting the number of columns and breaking
# lines if needed.
# Addionaly it expects you to pass a second argument which will be printed
# as an informational message.
# Last, but not least, a third command would instruct the function not to 
# executed the command.
print::command() {
    local SYMBOL=";"
    
    local COMMAND="$1"
    local COMMAND_DESCRIPTION="${2:-}"
    local COMMAND_MUST_BE_EXECUTED="${3:-1}"
    local MESSAGE=$(echo "${COMMAND}" | fold -w ${THEME_NUMBER_OF_COLUMNS} -s)
    local LINE=$(printf "%$((${THEME_NUMBER_OF_COLUMNS}))s" "$SYMBOL"| tr ' ;' "${THEME_SECTIOM_LINE_SYMBOL}")

    mecho "\n[dim]${LINE}[/]\n"
    mecho "[bold]\$[/] [bold][italic]${MESSAGE}[/]"

    if [[ $COMMAND_DESCRIPTION != "" ]] 
    then 
        print::info "${COMMAND_DESCRIPTION}"
    fi

    printf "\n"

    if [ $COMMAND_MUST_BE_EXECUTED -eq "0" ] 
    then 
        return
    fi

    local start=`date +%s`
    eval "${COMMAND}"
    local is_success=$?

    local end=`date +%s`
    local runtime=$((end-start))
    local ELAPSED=$(echo "$(($runtime % 60))s" | fold -w ${THEME_NUMBER_OF_COLUMNS} -s)

    if [ $is_success -eq "0" ] 
    then 
        print::success "It took \"${ELAPSED}\" to complete this job."
    else 
        print::error "It took \"${ELAPSED}\" to complete this job."
    fi
}
