function convertDate(dateInMs) {
    const date = new Date(dateInMs);
    const day = date.getUTCDate() > 9 ? date.getUTCDate() : `0${date.getUTCDate()}`;
    const monthInc = date.getMonth() + 1;
    const month = monthInc > 9 ? monthInc : `0${monthInc}`;
    const hours = date.getHours() > 9 ? date.getHours() : `0${date.getHours()}`;
    const minutes = date.getMinutes() > 9 ? date.getMinutes() : `0${date.getMinutes()}`;
    return `${date.getFullYear()}.${month}.${day} ${hours}:${minutes}`;
}

function getQueryParam(paramName) {
    let params = (new URL(document.location)).searchParams;
    return params.get(paramName)
}

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