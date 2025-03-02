let postToDelete = null; // Храним ID поста, который будем удалять

// Открытие модалки
function confirmDelete(postID) {
    postToDelete = postID;
    document.getElementById("deleteModal").style.display = "flex";
}

// Закрытие модалки
function closeModal() {
    postToDelete = null;
    document.getElementById("deleteModal").style.display = "none";
}

// Удаление поста
function deletePost() {
    if (!postToDelete) return;

    fetch(`/posts/${postToDelete}`, {
        method: "DELETE",
    })
    .then(response => {
        if (response.ok) {
            document.getElementById(`post-${postToDelete}`).remove(); // Удаляем из DOM
            closeModal();
        } else {
            alert("Ошибка при удалении поста");
        }
    })
    .catch(error => {
        console.error("Ошибка:", error);
    });
}
