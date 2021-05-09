function submitForm() {
    var f = document.getElementsByName('upload-form')[0];
    f.submit();
    f.reset();
}

function copyText() {
    const tmp = document.createElement('textarea');
    tmp.value = (document.getElementById("uploadcmd")).innerHTML;
    document.body.appendChild(tmp);
    tmp.select();
    document.execCommand('copy');
    document.body.removeChild(tmp);
}

function sortTable(idx, type) {
    var table, rows, switching, i, x, y, shouldSwitch, sortorder;
    table = document.getElementById("filetable");
    sortorder = document.getElementById("th_" + idx);
    switching = true;
    while (switching) {
        switching = false;
        rows = table.rows;
        for (i = 1; i < (rows.length - 1); i++) {
            shouldSwitch = false;
            x = rows[i].getElementsByTagName("TD")[idx];
            y = rows[i + 1].getElementsByTagName("TD")[idx];

            // comparison here
            // type 1:number 2:date
            if (type == 1) {
                if (Number(x.innerHTML) > Number(y.innerHTML) && sortorder.innerHTML > 0) {
                    shouldSwitch = true;
                    break;
                }
                
                if (Number(x.innerHTML) < Number(y.innerHTML) && sortorder.innerHTML < 0) {
                    shouldSwitch = true;
                    break;
                }
            } else if (type == 2){
                x = new Date(x.innerHTML);
                y = new Date(y.innerHTML);
                if (x > y && sortorder.innerHTML > 0) {
                    shouldSwitch = true;
                    break;
                }
                
                if (x < y && sortorder.innerHTML < 0) {
                    shouldSwitch = true;
                    break;
                }
            } else {
                if (x.innerHTML.toLowerCase() > y.innerHTML.toLowerCase() && sortorder.innerHTML > 0) {
                    shouldSwitch = true;
                    break;
                }

                if (x.innerHTML.toLowerCase() < y.innerHTML.toLowerCase() && sortorder.innerHTML < 0) {
                    shouldSwitch = true;
                    break;
                }
            }
        }
        if (shouldSwitch) {
            rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
            switching = true;
        }
    }
    sortorder.innerHTML = Number(sortorder.innerHTML) * -1;
}