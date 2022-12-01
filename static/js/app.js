function notify(msg, msgType) {
    notie.alert({
        type: msgType,
        text: msg
    })
}
function SuccessAlert(msg) {
    notie.alert({
    type: 1, // optional, default = 4, enum: [1, 2, 3, 4, 5, 'success', 'warning', 'error', 'info', 'neutral']
    text: msg, 
    // stay: Boolean, // optional, default = false
    // time: Number, // optional, default = 3, minimum = 1,
    position: 'top' // optional, default = 'top', enum: ['top', 'bottom']
    })
}
function ErrorAlert(msg) {
    notie.alert({
    type: 3, // optional, default = 4, enum: [1, 2, 3, 4, 5, 'success', 'warning', 'error', 'info', 'neutral']
    text: msg, 
    // stay: Boolean, // optional, default = false
    // time: Number, // optional, default = 3, minimum = 1,
    position: 'top' // optional, default = 'top', enum: ['top', 'bottom']
    })
}
function SweetErrorAlert(msg) {
    Swal.fire({
    icon: 'error',
    title: 'Oops...',
    showConfirmButton: false
    })
}

function Prompt() {
    let toast = function(c) {
        const {
            msg = "",
            icon = "",
            position = "top-end",

        } = c;
        const Toast = Swal.mixin({
            toast: true,
            title: msg, 
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
            toast.addEventListener('mouseenter', Swal.stopTimer)
            toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
        }
        let success = function(c) {
        const {
            msg = "Successfully",
            icon = "success",
            position = "center",

        } = c;
        Swal.fire({
            position: position,
            title: msg,
            icon: icon,
            showConfirmButton: false,
            timer: 1500
        })
    }
    let error = function(c) {
        const {
            msg = "Ooops ...",
            icon = "error",
            position = "center",

        } = c;
        Swal.fire({
            position: position,
            title: msg,
            icon: icon,
            showConfirmButton: false,
            timer: 1500
        })
    }

    let custom = async function(c) {
        const {
            msg = "",
            title = "",
            icon = "",
            showConfirmButton = true,
        } = c;

        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false, 
            focusConfirm: false,
            showCancelButton: true, 
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen()
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen()
                }
            }
        })

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                }else {
                    c.callback(false);
                }
            }else {
                c.callback(false);
            }
        }    
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom
    }
}


function CheckingDay(roomID, csrf) {
    let html = `
            <form action="/search-availability-json" method="post" id="check-availability-form" novalidate class="needs-validation">
                <div class="form-row">
                    <div class="col">
                    <div class="form-row" id="reservation-dates-modal">
                        <div class="col">
                            <input type="text" id="start" name="start" disabled required class="form-control" placeholder="Arrival">
                        </div>
                        <div class="col">
                            <input type="text" id="end" name="end" disabled required class="form-control" placeholder="Departure">
                        </div>
                    </div>
                    </div>
                </div>
            </form>
        `
    document.getElementById("check-availability-button").addEventListener("click", function() {
        attention.custom({
            msg: html, 
            title: "Choose your date",
            willOpen: () => {
                    const elem = document.getElementById("reservation-dates-modal");
                    const rp = new DateRangePicker(elem, {
                        format: 'yyyy-mm-dd',
                        showOnFocus: true,
                        minDate: new Date(),
                    })
            },
            didOpen: () => {
                    document.getElementById("start").removeAttribute('disabled');
                    document.getElementById("end").removeAttribute('disabled');
            },
            callback: (result) => {
                console.log("called back")

                let form = document.getElementById('check-availability-form');
                let formData = new FormData(form);
                formData.append("csrf_token", csrf);
                formData.append("room_id", roomID)

                fetch('/search-availability-json', {
                    method: "post",
                    body: formData,
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                icon: 'success',
                                showConfirmButton: false,
                                msg: '<p>Room is available!</p>'
                                    + '<p><a href="/book-room?id='
                                    + data.room_id
                                    + "&sd="
                                    + data.start_date
                                    + "&ed="
                                    + data.end_date     
                                    + '" class="btn btn-primary">Make Reservation Now</a></p>'
                                
                            })
                        }else {
                            attention.error({
                                msg: "Un Avalability"
                            })
                        }
                    })
            }
        })
     })
}