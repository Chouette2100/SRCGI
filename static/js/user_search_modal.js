// static/js/user_search_modal.js

/**
 * ユーザー検索モーダルを初期化する関数
 * @param {object} options - 初期化オプション
 * @param {string} options.modalID - モーダル要素のID (例: 'userSearchModal')
 * @param {string} options.openButtonID - モーダルを開くボタンのID (例: 'openUserSearchModal')
 * @param {string} options.IDPrefix - モーダル内の要素のIDプレフィックス (例: 'userSearch')
 * @param {function(string, string): void} options.onUserSelected - ユーザーが選択されたときに呼び出されるコールバック関数 (userno, username)
 */
function initUserSearchModal(options) {
    const { modalID, openButtonID, IDPrefix, onUserSelected } = options;

    const modal = document.getElementById(modalID);
    const openModalBtn = document.getElementById(openButtonID);
    const closeBtn = modal.querySelector(`[data-modal-close-target="${modalID}"]`); // data属性で閉じるボタンを特定
    const searchUserNameInput = document.getElementById(`${IDPrefix}SearchUserName`);
    const searchUserBtn = document.getElementById(`${IDPrefix}SearchUserBtn`);
    const searchResultsList = document.getElementById(`${IDPrefix}SearchResultsList`);
    const noResultsMessage = document.getElementById(`${IDPrefix}NoResultsMessage`);
    const selectUserAndRefreshBtn = document.getElementById(`${IDPrefix}SelectUserAndRefreshBtn`);

    let selectedUserNo = null;
    let selectedUserName = null;

    // ポップアップを開く
    if (openModalBtn) {
        openModalBtn.onclick = function() {
            modal.style.display = 'flex';
            searchUserNameInput.value = '';
            searchResultsList.innerHTML = '';
            noResultsMessage.style.display = 'none';
            selectUserAndRefreshBtn.disabled = true;
            selectedUserNo = null;
            selectedUserName = null;
        };
    }

    // ポップアップを閉じる (Xボタン)
    if (closeBtn) {
        closeBtn.onclick = function() {
            modal.style.display = 'none';
        };
    }

    // ポップアップを閉じる (モーダルの外側をクリック)
    window.addEventListener('click', function(event) {
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });

    // ユーザー名検索ボタン押下時の処理
    if (searchUserBtn) {
        searchUserBtn.onclick = function() {
            const nameQuery = searchUserNameInput.value;
            searchResultsList.innerHTML = '';
            noResultsMessage.style.display = 'none';
            selectUserAndRefreshBtn.disabled = true;
            selectedUserNo = null;
            selectedUserName = null;

            fetch(`/search-users?name=${encodeURIComponent(nameQuery)}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(users => {
                if (users.length === 0) {
                    noResultsMessage.style.display = 'block';
                } else {
                    users.forEach(user => {
                        const li = document.createElement('li');
                        li.textContent = `${user.username} (No: ${user.userno})`;
                        li.dataset.userno = user.userno;
                        li.dataset.username = user.username;
                        li.onclick = function() {
                            document.querySelectorAll(`#${IDPrefix}SearchResultsList li`).forEach(item => {
                                item.classList.remove('selected');
                            });
                            li.classList.add('selected');
                            selectedUserNo = user.userno;
                            selectedUserName = user.username;
                            selectUserAndRefreshBtn.disabled = false;
                        };
                        searchResultsList.appendChild(li);
                    });
                }
            })
            .catch(error => {
                console.error('Error searching users:', error);
                alert('ユーザー検索に失敗しました。');
            });
        };
    }

    // 「選択して決定」ボタン押下時の処理
    if (selectUserAndRefreshBtn) {
        selectUserAndRefreshBtn.onclick = function() {
            if (!selectedUserNo) {
                alert('ユーザーを選択してください。');
                return;
            }
            // コールバック関数を呼び出して、選択されたユーザー情報を親に渡す
            if (typeof onUserSelected === 'function') {
                onUserSelected(selectedUserNo, selectedUserName);
            }
            modal.style.display = 'none'; // ポップアップを閉じる
        };
    }
}
