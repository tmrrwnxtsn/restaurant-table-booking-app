function openModal(id, name) {
    const bookingForm = document.getElementById("makeBookingForm");
    bookingForm.setAttribute("action", `/restaurants/${id}/booked`);

    document.querySelector(".modal-header > h5").innerHTML = `Подтверждение брони в ресторане «${name}»`
    const arrayP = document.querySelectorAll(".row.g-3 > p")

    const peopleNumber = getQueryParam("people_number")

    arrayP[0].innerHTML = `Количество человек: ${peopleNumber}`
    document.querySelector("#people_number_input").value = peopleNumber

    const desiredDatetime = getQueryParam("desired_datetime")
    const dateObject = new Date(desiredDatetime)
    const convertedDate = convertDate(dateObject.getTime())

    arrayP[1].innerHTML = `Желаемое дата и время посещения: ${convertedDate}`
    document.querySelector("#desired_datetime_input").value = convertedDate
}

function getQueryParam(paramName) {
    let params = (new URL(document.location)).searchParams;
    return params.get(paramName)
}