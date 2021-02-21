//Used for messages on display on our messageboard
var MessageData = {
    TestString: "nameVal",
    TheMessages: new Array()
};
//Used to print username values on messages
var TheUser = {
    UserName: "",
    UserID: 0
};

var whatPage; //The page number we are on
var whatBoard; //which board we are currently on
var variableNameInt = 0; //A number to attach to variable IDs to make them unique

//This func intializes what number of comments we are loading
function initializeWhatPage(thePage){
    whatPage = thePage;
}
//This sets the username for TheUser
function setUsername(username){
    TheUser.UserName = username;
}
//This function sets what messageboard type we are on
function setBoardType(boardType){
    whatBoard = boardType;
}


//This sets the UserID for TheUser
function setUserID(userID){
    TheUser.UserID = userID;
}

//This puts data into our messageboard, usually called when the page refreshes
function messageDataInitialize(theData){
    /* Remove all Messages everytime this is called in order to populate the next page's
    data */
    var messageboardsectionDiv = document.getElementById("messageboardsectionDiv");
    messageboardsectionDiv.innerHTML = "";

    //Initialize variables above for messages
    var originalMessageArray = new Array(); //Declare array of message to go into
    var i;
    for (i = 0; i < theData.length; i++) {
        originalMessageArray.push(theData[i]);
        //console.log("Here is our original message data: " + originalMessageArray[i].TheMessage);
    }
    //Display current page value
    changePageNumber(whatPage);
    /*
        Add messages to the main board, (if they are not child messages). If they are child messages, 
        recursivley create 'messageSectionDiv' as replies
    */
   recursiveMessageAdd(originalMessageArray, messageboardsectionDiv);
}

//This function recursivley creates a space for messages and their replies
//It is called multpile times recursivley for each reply to a reply and so on
function recursiveMessageAdd(currentArray, parentDiv){
    var i = 0;
    /*For each message found in the given messageArray, attach it to the parent div
    ,(then, if it has messages of it's own, call this function again for that message
    that is ALSO a parent). */
    var messageArray = new Array(); //Declare array of message to go into
    var j = 0;
    for (j = 0; j < currentArray.length; j++) {
        messageArray.push(currentArray[j]);
    }
    for (i = 0; i < messageArray.length; i++){
        //Increment variable name to make our names unique
        variableNameInt = variableNameInt + 1;
        //console.log("Currently making message space for this message: " + messageArray[i].TheMessage);
        createMessageDivs(variableNameInt, messageArray[i], parentDiv);
    }
}
/* 
    This is called in the recursiveMessageAdd function; it is primarily used
    for preloading messages on new page loads and adding replies to those
    messages recursively. We have another function for adding 
    new messages to our parent div
*/
function createMessageDivs(variableNameIntCurrently, currentMessage, currentParent){
    /* Create elements to attach to the parent */
    var idNamer = ""; //Declare ID Namer for the various ID names we're giving
    var varNameInt = Number(variableNameIntCurrently); //Used for number designation in naming
    /* MESSAGE SECTION */
    //messageSectionDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "messageSectionDiv" + varNameInt.toString();
    var messageSectionDiv = document.createElement("div");
    messageSectionDiv.setAttribute("id", idNamer);
    messageSectionDiv.setAttribute("class", "messageSectionDiv");
    //messageTextDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "messageTextDiv" + varNameInt.toString();
    var messageTextDiv = document.createElement("div");
    messageTextDiv.setAttribute("id", idNamer);
    messageTextDiv.setAttribute("class", "messageTextDiv");
    //textP
    idNamer = ""; //Set idNamer to nil
    idNamer = "textP" + varNameInt.toString();
    var textP = document.createElement("p");
    textP.setAttribute("id", idNamer);
    textP.setAttribute("class", "textP");
    textP.innerHTML = currentMessage.TheMessage; //Set Message value
    //Attach Elements to each other
    messageTextDiv.appendChild(textP);
    messageSectionDiv.appendChild(messageTextDiv);

    /* REPLY BUTTON SECTION */
    //replyButtonSpacingDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyButtonSpacingDiv" + varNameInt.toString();
    var replyButtonSpacingDiv = document.createElement("div");
    replyButtonSpacingDiv.setAttribute("id", idNamer);
    replyButtonSpacingDiv.setAttribute("class", "replyButtonSpacingDiv");
    //replyButtonDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyButtonDiv" + varNameInt.toString();
    var replyButtonDiv = document.createElement("div");
    replyButtonDiv.setAttribute("id", idNamer);
    replyButtonDiv.setAttribute("class", "replyButtonDiv");
    //buttonTextDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "buttonTextDiv" + varNameInt.toString();
    var buttonTextDiv = document.createElement("div");
    buttonTextDiv.setAttribute("id", idNamer);
    buttonTextDiv.setAttribute("class", "buttonTextDiv");
    //replyNumP
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyNumP" + varNameInt.toString();
    var replyNumP = document.createElement("p");
    replyNumP.setAttribute("id", idNamer);
    replyNumP.setAttribute("class", "replyNumP");
    var numOReplies = Number(currentMessage.RepliesAmount);
    replyNumP.innerHTML = numOReplies.toString() + " Replies";
    //buttonIconDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "buttonIconDiv" + varNameInt.toString();
    var buttonIconDiv = document.createElement("div");
    buttonIconDiv.setAttribute("id", idNamer);
    buttonIconDiv.setAttribute("class", "buttonIconDiv");
    //Attach Elements to each other
    buttonTextDiv.appendChild(replyNumP);
    replyButtonDiv.appendChild(buttonTextDiv);
    replyButtonDiv.appendChild(buttonIconDiv);
    replyButtonSpacingDiv.appendChild(replyButtonDiv);
    messageSectionDiv.appendChild(replyButtonSpacingDiv);

    /* REPLIES SECTION */
    //repliesSectionDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "repliesSectionDiv" + varNameInt.toString();
    var repliesSectionDiv = document.createElement("div");
    repliesSectionDiv.setAttribute("id", idNamer);
    repliesSectionDiv.setAttribute("class", "repliesSectionDiv");
    repliesSectionDiv.style.display = "none";
    //replyListingDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyListingDiv" + varNameInt.toString();
    var replyListingDiv = document.createElement("div");
    replyListingDiv.setAttribute("id", idNamer);
    replyListingDiv.setAttribute("class", "replyListingDiv");
    //replyInputDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyInputDiv" + varNameInt.toString();
    var replyInputDiv = document.createElement("div");
    replyInputDiv.setAttribute("id", idNamer);
    replyInputDiv.setAttribute("class", "replyInputDiv");
    //textareaReply
    idNamer = ""; //Set idNamer to nil
    idNamer = "textareaReply" + varNameInt.toString();
    var textareaReply = document.createElement("textarea");
    textareaReply.setAttribute("id", idNamer);
    textareaReply.setAttribute("class", "textareaReply");
    textareaReply.setAttribute("name", "textareaReply");
    textareaReply.setAttribute("placeholder", "Reply to this comment...");
    //replyButtonSendMsg
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyButtonSendMsg" + varNameInt.toString();
    var replyButtonSendMsg = document.createElement("button");
    replyButtonSendMsg.setAttribute("id", idNamer);
    replyButtonSendMsg.setAttribute("class", "replyButtonSendMsg");
    replyButtonSendMsg.innerHTML = "Reply";
    //Attach elements to each other
    replyInputDiv.appendChild(textareaReply);
    replyInputDiv.appendChild(replyButtonSendMsg);
    repliesSectionDiv.appendChild(replyInputDiv);
    repliesSectionDiv.appendChild(replyListingDiv);
    messageSectionDiv.appendChild(repliesSectionDiv);

    /* Add onclick events for certain elements and 
    events to send data when it is entered */
    buttonIconDiv.addEventListener("click", function(){
        //Open replies if repliesSectionDiv is hidden
        if (repliesSectionDiv.style.display === "none"){
            //Replies are hidden, show them!
            repliesSectionDiv.style.display = "flex";
            buttonIconDiv.style.backgroundImage = 'url(static/images/svg/uparrow.svg)'; //Set Image
        } else {
            //Replies are showing, hide them!
            console.log("DEBUG: We are hiding the results for this buttonIcon: " + buttonIconDiv.id);
            repliesSectionDiv.style.display = "none";
            buttonIconDiv.style.backgroundImage = 'url(static/images/svg/downarrow.svg)'; //Set Image
        }
    }); //Event for when reply section needs to bee hidden/seen

    replyButtonSendMsg.addEventListener("click", function(){
        //Initial Check to see if anything is entered in the reply section
        if (textareaReply.value === ""){
            console.log("DEBUG: No response entered in the reply field");
        } else {
            /* Reply has been created; attach it to this parent in the backend, 
            then attach it as a div, recursivley calling this function */
            var TheParentMessage = {
                MessageID: Number(currentMessage.MessageID),
                UserID: Number(currentMessage.UserID),
                Messages: currentMessage.Messages,
                IsChild: currentMessage.IsChild,
                HasChildren: currentMessage.HasChildren,
                ParentMessageID: currentMessage.ParentMessageID,
                UberParentID: currentMessage.UberParentID,
                Order: currentMessage.Order,
                RepliesAmount: currentMessage.RepliesAmount,
                TheMessage: currentMessage.TheMessage,
                DateCreated: currentMessage.DateCreated,
                LastUpdated: currentMessage.LastUpdated
            }; //Declare Message JSON for parent,(which will be currentMessage)
            var TheChildMessage = {
                MessageID: 0,
                UserID: 0,
                Messages: new Array(),
                IsChild: true,
                HasChildren: false,
                ParentMessageID: Number(currentMessage.MessageID),
                UberParentID: currentMessage.UberParentID,
                Order: 0,
                RepliesAmount: 0,
                TheMessage: String(textareaReply.value),
                DateCreated: "",
                LastUpdated: ""
            }; //Declare Message JSON for parent,(which will be a message we make with reply to parent)
            var MessageReply = {
                ParentMessage: TheParentMessage,
                ChildMessage: TheChildMessage,
                CurrentPage: whatPage
            };
            //Reply with Ajax
            var jsonString = JSON.stringify(MessageReply); //Stringify Data
            //Send Request to change page
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/messageReplyAjax', true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.addEventListener('readystatechange', function(){
                if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                    var item = xhr.responseText;
                    var ReturnData = JSON.parse(item);
                    if (ReturnData.SuccessInt == 0){
                        /* Inserted reply correctly; we can now update this page with the new reply,
                        now that it is loaded onto our server and backend databases */
                        variableNameInt = variableNameInt + 1; //Increment for unique ID
                        var replyMessageDiv = createMessageReply(variableNameInt, ReturnData.CreatedMessage);
                        //Prepend this message BEHIND the reply section which should stay first in queue
                        replyListingDiv.prepend(replyMessageDiv);
                        //repliesSectionDiv.prepend(replyMessageDiv);
                        /* clear this text area's values */
                        textareaReply.value = ""; //Clear any values put in
                        textareaReply.placeholder = ReturnData.SuccessMsg; //Display Message
                        /* Update the amount of replies to this message */
                        numOReplies = numOReplies + 1;
                        replyNumP.innerHTML = numOReplies.toString() + " Replies";
                    } else {
                        /* Page Load unsuccessfull. Informing User in goInput */
                        textareaReply.value = ""; //Clear any values put in
                        textareaReply.innerHTML = ""; //Clear any value put in there
                        textareaReply.placeholder = ReturnData.SuccessMsg; //Display Message
                    }
                }
            });
            xhr.send(jsonString);
        }
    }); //Event for when we need to make a reply with Ajax

    //Attach the message section div to this parent
    currentParent.appendChild(messageSectionDiv);

    /* 
        Recursivley add messages to this parent is there are replies to this
        message/reply
    */
    if (currentMessage.HasChildren === true){
        //Need the array backwards to show the most recent messages
        var reversedArray = new Array();
        reversedArray = currentMessage.Messages;
        reversedArray.reverse();
        reversedArray.length
        //Start creating divs for the children
        var g = 0;
        for (g = 0; g < reversedArray.length; g++){
            /*
            console.log("Adding replied messages to this message: " + reversedArray[g].TheMessage);
            console.log("Here is the " + g + " message we'll add for: " + reversedArray[g].TheMessage);
            */
            var simpleArray = new Array();
            simpleArray.push(reversedArray[g]);
            recursiveMessageAdd(simpleArray, repliesSectionDiv);
        }
    } else {
        //console.log("DEBUG: This has no more parent left");
    }
}

/* This function is called within 'createMessageDivs'; it immediatley creates a 
div to attach to the message it was replying to */
function createMessageReply(variableNameIntCurrently, currentMessage){
    /* Create elements to attach to the parent */
    var idNamer = ""; //Declare ID Namer for the various ID names we're giving
    var varNameInt = Number(variableNameIntCurrently); //Used for number designation in naming
    /* MESSAGE SECTION */
    //messageSectionDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "messageSectionDiv" + varNameInt.toString();
    var messageSectionDiv = document.createElement("div");
    messageSectionDiv.setAttribute("id", idNamer);
    messageSectionDiv.setAttribute("class", "messageSectionDiv");
    //messageTextDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "messageTextDiv" + varNameInt.toString();
    var messageTextDiv = document.createElement("div");
    messageTextDiv.setAttribute("id", idNamer);
    messageTextDiv.setAttribute("class", "messageTextDiv");
    //textP
    idNamer = ""; //Set idNamer to nil
    idNamer = "textP" + varNameInt.toString();
    var textP = document.createElement("p");
    textP.setAttribute("id", idNamer);
    textP.setAttribute("class", "textP");
    textP.innerHTML = currentMessage.TheMessage; //Set Message value
    //Attach Elements to each other
    messageTextDiv.appendChild(textP);
    messageSectionDiv.appendChild(messageTextDiv);

    /* REPLY BUTTON SECTION */
    //replyButtonSpacingDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyButtonSpacingDiv" + varNameInt.toString();
    var replyButtonSpacingDiv = document.createElement("div");
    replyButtonSpacingDiv.setAttribute("id", idNamer);
    replyButtonSpacingDiv.setAttribute("class", "replyButtonSpacingDiv");
    //replyButtonDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyButtonDiv" + varNameInt.toString();
    var replyButtonDiv = document.createElement("div");
    replyButtonDiv.setAttribute("id", idNamer);
    replyButtonDiv.setAttribute("class", "replyButtonDiv");
    //buttonTextDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "buttonTextDiv" + varNameInt.toString();
    var buttonTextDiv = document.createElement("div");
    buttonTextDiv.setAttribute("id", idNamer);
    buttonTextDiv.setAttribute("class", "buttonTextDiv");
    //replyNumP
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyNumP" + varNameInt.toString();
    var replyNumP = document.createElement("p");
    replyNumP.setAttribute("id", idNamer);
    replyNumP.setAttribute("class", "replyNumP");
    var numOReplies = Number(currentMessage.RepliesAmount);
    replyNumP.innerHTML = numOReplies.toString() + " Replies";
    //buttonIconDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "buttonIconDiv" + varNameInt.toString();
    var buttonIconDiv = document.createElement("div");
    buttonIconDiv.setAttribute("id", idNamer);
    buttonIconDiv.setAttribute("class", "buttonIconDiv");
    //Attach Elements to each other
    buttonTextDiv.appendChild(replyNumP);
    replyButtonDiv.appendChild(buttonTextDiv);
    replyButtonDiv.appendChild(buttonIconDiv);
    replyButtonSpacingDiv.appendChild(replyButtonDiv);
    messageSectionDiv.appendChild(replyButtonSpacingDiv);

    /* REPLIES SECTION */
    //repliesSectionDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "repliesSectionDiv" + varNameInt.toString();
    var repliesSectionDiv = document.createElement("div");
    repliesSectionDiv.setAttribute("id", idNamer);
    repliesSectionDiv.setAttribute("class", "repliesSectionDiv");
    repliesSectionDiv.style.display = "none";
    //replyListingDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyListingDiv" + varNameInt.toString();
    var replyListingDiv = document.createElement("div");
    replyListingDiv.setAttribute("id", idNamer);
    replyListingDiv.setAttribute("class", "replyListingDiv");
    //replyInputDiv
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyInputDiv" + varNameInt.toString();
    var replyInputDiv = document.createElement("div");
    replyInputDiv.setAttribute("id", idNamer);
    replyInputDiv.setAttribute("class", "replyInputDiv");
    //textareaReply
    idNamer = ""; //Set idNamer to nil
    idNamer = "textareaReply" + varNameInt.toString();
    var textareaReply = document.createElement("textarea");
    textareaReply.setAttribute("id", idNamer);
    textareaReply.setAttribute("class", "textareaReply");
    textareaReply.setAttribute("name", "textareaReply");
    textareaReply.setAttribute("placeholder", "Reply to this comment...");
    //replyButtonSendMsg
    idNamer = ""; //Set idNamer to nil
    idNamer = "replyButtonSendMsg" + varNameInt.toString();
    var replyButtonSendMsg = document.createElement("button");
    replyButtonSendMsg.setAttribute("id", idNamer);
    replyButtonSendMsg.setAttribute("class", "replyButtonSendMsg");
    replyButtonSendMsg.innerHTML = "Reply";
    //Attach elements to each other
    replyInputDiv.appendChild(textareaReply);
    replyInputDiv.appendChild(replyButtonSendMsg);
    repliesSectionDiv.appendChild(replyInputDiv);
    repliesSectionDiv.appendChild(replyListingDiv);
    messageSectionDiv.appendChild(repliesSectionDiv);

    /* Add onclick events for certain elements and 
    events to send data when it is entered */
    buttonIconDiv.addEventListener("click", function(){
        //Open replies if repliesSectionDiv is hidden
        if (repliesSectionDiv.style.display === "none"){
            //Replies are hidden, show them!
            repliesSectionDiv.style.display = "flex";
            buttonIconDiv.style.backgroundImage = 'url(static/images/svg/uparrow.svg)'; //Set Image
        } else {
            //Replies are showing, hide them!
            repliesSectionDiv.style.display = "none";
            buttonIconDiv.style.backgroundImage = 'url(static/images/svg/downarrow.svg)'; //Set Image
        }
    }); //Event for when reply section needs to bee hidden/seen

    replyButtonSendMsg.addEventListener("click", function(){
        //Initial Check to see if anything is entered in the reply section
        if (textareaReply.value === ""){
            console.log("DEBUG: No response entered in the reply field");
        } else {
            /* Reply has been created; attach it to this parent in the backend, 
            then attach it as a div, recursivley calling this function */
            var TheParentMessage = {
                MessageID: Number(currentMessage.MessageID),
                UserID: Number(currentMessage.UserID),
                Messages: currentMessage.Messages,
                IsChild: currentMessage.IsChild,
                HasChildren: currentMessage.HasChildren,
                ParentMessageID: currentMessage.ParentMessageID,
                UberParentID: currentMessage.UberParentID,
                Order: currentMessage.Order,
                RepliesAmount: currentMessage.RepliesAmount,
                TheMessage: currentMessage.TheMessage,
                DateCreated: currentMessage.DateCreated,
                LastUpdated: currentMessage.LastUpdated
            }; //Declare Message JSON for parent,(which will be currentMessage)
            var TheChildMessage = {
                MessageID: 0,
                UserID: 0,
                Messages: new Array(),
                IsChild: true,
                HasChildren: false,
                ParentMessageID: Number(currentMessage.MessageID),
                UberParentID: currentMessage.UberParentID,
                Order: 0,
                RepliesAmount: 0,
                TheMessage: String(textareaReply.value),
                DateCreated: "",
                LastUpdated: ""
            }; //Declare Message JSON for parent,(which will be a message we make with reply to parent)
            var MessageReply = {
                ParentMessage: TheParentMessage,
                ChildMessage: TheChildMessage,
                CurrentPage: whatPage
            };
            //Reply with Ajax
            var jsonString = JSON.stringify(MessageReply); //Stringify Data
            //Send Request to change page
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/messageReplyAjax', true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.addEventListener('readystatechange', function(){
                if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                    var item = xhr.responseText;
                    var ReturnData = JSON.parse(item);
                    if (ReturnData.SuccessInt == 0){
                        /* Inserted reply correctly; we can now update this page with the new reply,
                        now that it is loaded onto our server and backend databases */
                        variableNameInt = variableNameInt + 1; //Increment for unique ID
                        var replyMessageDiv = createMessageReply(variableNameInt, ReturnData.CreatedMessage);
                        //Prepend this message BEHIND the reply section which should stay first in queue
                        replyListingDiv.prepend(replyMessageDiv);
                        //repliesSectionDiv.prepend(replyMessageDiv);
                        /* clear this text area's values */
                        textareaReply.value = ""; //Clear any values put in
                        textareaReply.placeholder = ReturnData.SuccessMsg; //Display Message
                        /* Update the amount of replies to this message */
                        numOReplies = numOReplies + 1;
                        replyNumP.innerHTML = numOReplies.toString() + " Replies";
                    } else {
                        /* Page Load unsuccessfull. Informing User in goInput */
                        textareaReply.value = ""; //Clear any values put in
                        textareaReply.innerHTML = ""; //Clear any value put in there
                        textareaReply.placeholder = ReturnData.SuccessMsg; //Display Message
                    }
                }
            });
            xhr.send(jsonString);
        }
    }); //Event for when we need to make a reply with Ajax

    //return the final div to the parent
    return messageSectionDiv;
}

//Add event to check if this Message board exists on this page; if so, set value to 1 on load of page
window.addEventListener('DOMContentLoaded', function(){
    //Declare variables that should be on this page
    var whatPageP = document.getElementById("whatPageP"); //Lists what page it is
    var goToButton = document.getElementById("goToButton"); //Button for initializing page number
    var leftButtonDiv = document.getElementById("leftButtonDiv"); //Left Button Page Navigation
    var rightButtonDiv = document.getElementById("rightButtonDiv"); //Right Button Page Navigation

    //Set page number on initial page loading
    if (whatPageP === null){
        console.log("DEBUG: This page does not have the messageboard");
    } else {
        changePageNumber(whatPage);
    }

    //Set the 'onClick' event for this goToButton
    goToButton.addEventListener("click", function(){
        //Declare the input field
        var goInput = document.getElementById("goInput");
        var thePage = Number(goInput.value);
        var PageData = {
            ThePage: thePage,
        }
        var jsonString = JSON.stringify(PageData); //Stringify Data
        //Send Request to change page
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/evaluateTenResults', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                var ReturnMessage = JSON.parse(item);
                if (ReturnMessage.SuccOrFail == 0){
                    /* Page Load successful, change data for new page loaded */
                    whatPage = thePage; //Set value of page for future page use.
                    var MessageViewData = {
                        TestString: ReturnMessage.ResultMsg,
                        TheMessages: ReturnMessage.Messages
                    }
                    //location.reload();
                    messageDataInitialize(ReturnMessage.Messages); //Send Messages to messageDataInitialize
                } else {
                    /* Page Load unsuccessfull. Informing User in goInput */
                    document.getElementById("goInput").value = ''; //Clear any values put in
                    document.getElementById("goInput").placeholder = ReturnMessage.ResultMsg; //Display Message
                }
            }
        });
        xhr.send(jsonString);
    });

    //Set the 'onClick' event for this leftButtonDiv
    leftButtonDiv.addEventListener("click", function(){
        var thePage = whatPage -1;
        var PageData = {
            ThePage: thePage,
        }
        var jsonString = JSON.stringify(PageData); //Stringify Data
        //Send Request to change page
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/evaluateTenResults', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                var ReturnMessage = JSON.parse(item);
                if (ReturnMessage.SuccOrFail == 0){
                    /* Page Load successful, change data for new page loaded */
                    whatPage = thePage; //Set value of page for future page use.
                    var MessageViewData = {
                        TestString: ReturnMessage.ResultMsg,
                        TheMessages: ReturnMessage.Messages
                    }
                    //location.reload();
                    messageDataInitialize(ReturnMessage.Messages); //Send Messages to messageDataInitialize
                } else {
                    /* Page Load unsuccessfull. Informing User in goInput */
                    document.getElementById("goInput").value = ''; //Clear any values put in
                    document.getElementById("goInput").placeholder = ReturnMessage.ResultMsg; //Display Message
                }
            }
        });
        xhr.send(jsonString);
    });
    //Set the 'onClick' event for this rightButtonDiv
    rightButtonDiv.addEventListener("click", function(){
        var thePage = whatPage + 1;
        var PageData = {
            ThePage: thePage,
        }
        var jsonString = JSON.stringify(PageData); //Stringify Data
        //Send Request to change page
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/evaluateTenResults', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                var ReturnMessage = JSON.parse(item);
                if (ReturnMessage.SuccOrFail == 0){
                    /* Page Load successful, change data for new page loaded */
                    whatPage = thePage; //Set value of page for future page use.
                    var MessageViewData = {
                        TestString: ReturnMessage.ResultMsg,
                        TheMessages: ReturnMessage.Messages
                    }
                    messageDataInitialize(ReturnMessage.Messages); //Send Messages to messageDataInitialize
                    //location.reload();
                } else {
                    /* Page Load unsuccessfull. Informing User in goInput */
                    document.getElementById("goInput").value = ''; //Clear any values put in
                    document.getElementById("goInput").placeholder = ReturnMessage.ResultMsg; //Display Message
                }
            }
        });
        xhr.send(jsonString);
    });
});

//Change the Page Number showing to what page it currently is
function changePageNumber(thePageCurrently){
    var whatPageP = document.getElementById("whatPageP");
    //Change the innerHTML of whatPageP to include the page number
    var whatPageText = "Page ";
    var pageNum = Number(thePageCurrently);
    var pageToString = pageNum.toString();
    whatPageText = whatPageText + pageToString;
    whatPageP.innerHTML = whatPageText;
    //Change the value of goInput to reflect the change
    document.getElementById("goInput").value = thePageCurrently; //Clear any values put in
    document.getElementById("goInput").placeholder = thePageCurrently; //Display Message
}
//Creates an original comment when the 'Add Comment' button is clicked.
function orignalCommentMaker(){
    console.log("DEBUG: Submitting an original comment.");
    var textareaComment = document.getElementById("textareaComment");
    var OriginalMessage = {
        TheMessage: String(textareaComment.value),
        PosterName: String(TheUser.UserName),
        WhatBoard: String(whatBoard)
    };
    //Original Message with Ajax
    var jsonString = JSON.stringify(OriginalMessage); //Stringify Data
    //Send Request to change page
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/messageOriginalAjax', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var DataReturn = JSON.parse(item);
            if (DataReturn.SuccessInt == 0){
                /* Inserted new, original message correctly; we can now update this page with the new message,
                now that it is loaded onto our server and backend databases */
                location.reload();
            } else {
                /* Page Load unsuccessfull. Informing User in goInput */
                textareaComment.value = ""; //Clear any values put in
                textareaComment.innerHTML = ""; //Clear any value put in theree
                textareaComment.placeholder = DataReturn.SuccessMsg; //Display error message
            }
        }
    });
    xhr.send(jsonString);
}

/* Test area stuff */

function testShowData(theData){
    console.log("Here is the data: " + theData);
    console.log("Here is but one bit of data: " + theData[0].TheMessage);
    for (var i = 0; i < theData.length; i++){
        console.log("Here is our data for " + i + ": " + theData[i].TheMessage);
    }
}