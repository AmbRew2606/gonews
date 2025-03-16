let postToDelete = null; 

function confirmDelete(postID) {
    postToDelete = postID;
    document.getElementById("deleteModal").style.display = "flex";
}

function closeModal() {
    postToDelete = null;
    document.getElementById("deleteModal").style.display = "none";
}


function deletePost() {
    if (!postToDelete) return;

    fetch(`/posts/${postToDelete}`, {
        method: "DELETE",
    })
    .then(response => {
        if (response.ok) {
            document.getElementById(`post-${postToDelete}`).remove();
            closeModal();
        } else {
            alert("Ошибка при удалении поста");
        }
    })
    .catch(error => {
        console.error("Ошибка:", error);
    });
}
