/* placeholder file for JavaScript */

/* format date "YYYY-MM-DD HH:MM:SS" */
function format_date(date_str) {
    date_str = String(date_str)
    date_str = date_str.split("T").join(" ");
    date_str = date_str.split("+")[0];
    return date_str;
}

function select_filter(option_str) {
    let elm = document.getElementById("select");
    if (option_str == null)
        elm.value = "all";
    else
        elm.value = option_str;
}