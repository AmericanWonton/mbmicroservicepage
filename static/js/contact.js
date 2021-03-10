function sendUpdate(){
    //Declare variables
    var YourNameInput = document.getElementById("YourNameInput");
    var YourEmailInput = document.getElementById("YourEmailInput");
    var YourMessageInput = document.getElementById("YourMessageInput");
    var informTxtP = document.getElementById("informTxtP");

    var MessageInfo = {
        YourNameInput: String(YourNameInput.value),
        YourEmailInput: Number(YourEmailInput.value),
        YourMessageInput: String(YourMessageInput.value)
    };
    var jsonString = JSON.stringify(MessageInfo); //Stringify Data
    //Send Request to user message update page
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/genericUserContact', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var ReturnData = JSON.parse(item);
            if (ReturnData.SuccessInt == 0){
                informTxtP.innerHTML = "Email sent. I'll write back soon!";
            } else {
                informTxtP.innerHTML = "Sorry! I couldn't get your messsage.... :(";
            }
        }
    });
    xhr.send(jsonString);
}