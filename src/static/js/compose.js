function openComposeModal() {
    var composeModal = document.getElementById('composeModal');
    composeModal.modal();
}

function attachComposeModal() {
    var composeAnnouncementBtn = document.getElementById('compose-announcements-btn');
    composeAnnouncementBtn.addEventListener('click', openComposeModal);
}

$(document).ready(function(){
    $('.modal').modal();
});
attachComposeModal();
