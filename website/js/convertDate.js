function convertDate(dateInMs) {
    const date = new Date(dateInMs);
    const day = date.getUTCDate() > 9 ? date.getUTCDate() : `0${date.getUTCDate()}`;
    const monthInc = date.getMonth() + 1
    const month = monthInc > 9 ? monthInc : `0${monthInc}`;
    const hours = date.getHours() > 9 ? date.getHours() : `0${date.getHours()}`;
    const minutes = date.getMinutes() > 9 ? date.getMinutes() : `0${date.getMinutes()}`;
    return `${day}.${month}.${date.getFullYear()} ${hours}:${minutes}`;
}