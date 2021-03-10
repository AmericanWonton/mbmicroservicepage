//Open the correctDivs when clicked
function documentDivDisplay(whichDiv){
    switch(whichDiv){
        case 1:
            //Display or not display 1st Div
            var theDiv = document.getElementById("bodOpen1");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 2:
            //Display or not display 2nd Div
            var theDiv = document.getElementById("bodOpen2");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 3:
            //Display or not display 3rd Div
            var theDiv = document.getElementById("bodOpen3");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 4:
            //Display or not display 4th Div
            var theDiv = document.getElementById("bodOpen4");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 5:
            //Display or not display 5th Div
            var theDiv = document.getElementById("bodOpen5");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 6:
            //Display or not display 6th Div
            var theDiv = document.getElementById("bodOpen6");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        default:
            console.log("Error, incorrect div was opened.");
            break;
    }
}