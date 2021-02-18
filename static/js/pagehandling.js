
//Used to control which link to send our user to
function navigateHeader(whichLink) {
    switch (whichLink) {
        case 1:
            //Go to ContactDev
            window.location.replace("/contact");
            break;
        case 2:
            //Go to Documentation
            window.location.replace("/documentation");
            break;
        case 3:
            //Go to hotdog Messageboard
            window.location.replace("/hotdogMB");
            break;
        case 4:
            //Go to hamburger messageboard
            window.location.replace("/hamburgerMB");
            break;
        case 5:
            //Go to Index
            window.location.replace("/");
            break;
        default:
            console.log("Error, wrong whichLink entered: " + whichLink);
            break;
    }
}