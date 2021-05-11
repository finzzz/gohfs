baseURL = document.getElementById("baseURL").innerHTML
zipPath = document.getElementById("zipPath").innerHTML
sha1Path = document.getElementById("sha1Path").innerHTML
init()

function init() {
    document.getElementById("up_curl").innerHTML = "curl -F 'file=@uploadthis.txt' " + baseURL

    renderItems("linksvg", "static/icons/link.svg")
    renderItems("zipsvg", "static/icons/zip.svg")
    renderItems("qrsvg", "static/icons/qr.svg")
    renderItems("hashsvg", "static/icons/hash.svg")
    renderItems("termsvg", "static/icons/term.svg")

    setPath(zipPath, "ziplink")
    setPath(sha1Path, "sha1link")
}

function setPath(path, cls) {
    link = document.getElementsByClassName(cls)

    Array.prototype.forEach.call(link, function(slide, index) {
        link.item(index).href = path + location.pathname + link.item(index).name
    });
}

function renderItems(cls, url) {
    webPath = document.getElementById("webPath")
    link = document.getElementsByClassName(cls)

    Array.prototype.forEach.call(link, function(slide, index) {
        link.item(index).src = webPath.innerHTML + "/" + url
    });
}

function submitForm() {
    f = document.getElementsByName("upload-form")[0]
    f.submit()
    f.reset()
}

function copyTextAsURL(text) {
    tmp = document.createElement("textarea")
    tmp.value = baseURL + location.pathname + text
    copy(tmp)
}

function copyTextById(id) {
    tmp = document.createElement("textarea")
    tmp.value = (document.getElementById(id)).innerHTML
    copy(tmp)
}

function copy(val) {
    document.body.appendChild(val)
    val.select()
    document.execCommand("copy")
    document.body.removeChild(val)
}

function getCellIdx(id) {
    if (id == "size" || id == "modtime") {
        return document.getElementById("th_"+id).cellIndex + 1
    }

    return document.getElementById("th_"+id).cellIndex
}

function compare(x, y, sortorder) {
    if ((x > y && sortorder > 0) || (x < y && sortorder < 0)) {
        return true
    }

    return false
}

function sortTable(id, type) {
    var rows, i, x, y, shouldSwitch
    table = document.getElementById("filetable")
    header = document.getElementById("th_" + id)
    sortorder = header.getAttribute("data-sortorder")

    cellIdx = getCellIdx(id)
    switching = true
    while (switching) {
        switching = false
        rows = table.rows
        for (i = 1; i < (rows.length - 1); i++) {
            shouldSwitch = false
            x = rows[i].getElementsByTagName("TD")[cellIdx]
            y = rows[i + 1].getElementsByTagName("TD")[cellIdx]

            if (type == "num") {
                if (compare(Number(x.innerHTML), Number(y.innerHTML), sortorder)) {
                    shouldSwitch = true
                    break
                }
            } else if (type == "date"){
                x = new Date(x.innerHTML)
                y = new Date(y.innerHTML)
                if (compare(x, y, sortorder)) {
                    shouldSwitch = true
                    break
                }
            } else {
                if (compare(x.innerHTML.toLowerCase(), y.innerHTML.toLowerCase(), sortorder)) {
                    shouldSwitch = true
                    break
                }
            }
        }

        if (shouldSwitch) {
            rows[i].parentNode.insertBefore(rows[i + 1], rows[i])
            switching = true
        }
    }

    header.setAttribute("data-sortorder", Number(sortorder) * -1)
}

function showModal(id) {
    modal = document.getElementById(id)
    modal.style.display = "flex"
}

function hideModal(id) {
    modal = document.getElementById(id)
    modal.style.display = "none"

    document.getElementById("qrcode").innerHTML = ""
    document.getElementById("zipqrcode").innerHTML = ""
}

function showQR(url) {
    qrdiv = document.getElementById("qrcode")
    qrcode = new QRCode(qrdiv, {width: 256, height: 256, margin: "auto"})
    link = baseURL + location.pathname + url
    document.getElementById("caption").innerHTML = "RAW : " + link

    // zip
    zipqrdiv = document.getElementById("zipqrcode")
    zipqrcode = new QRCode(zipqrdiv, {width: 256, height: 256, margin: "auto"})
    ziplink = baseURL + zipPath + location.pathname + url
    document.getElementById("zipcaption").innerHTML = "ZIP : " + ziplink

    showModal("QRModal")

    qrcode.makeCode(link)
    zipqrcode.makeCode(ziplink)
}