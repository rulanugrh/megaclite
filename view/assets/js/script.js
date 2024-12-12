document.addEventListener('DOMContentLoaded', function() {
    // create new variabel for catch element by id
    const sidebar = document.getElementById('sidebar');
    const hamburgerButton = document.getElementById('hamburgerButton');
    const hamburgerIcon = document.getElementById('hamburgerIcon');
    const closeButton = document.getElementById('closeButton');
    const mainContent = document.getElementById('mainContent');

    // variabel for get all items in navbar
    const sidebarItems = document.querySelectorAll('#mail a');
    
    // function to close navbar
    function closeSidebar() {
        sidebar.classList.add('-translate-x-full');
        hamburgerIcon.style.display = 'block'; 

        // condition to style close button in large and small screen
        if (window.innerWidth >= 1024) {
            closeButton.style.display = 'block';
        } else { // Small screen
            closeButton.style.display = 'none';
        }
        mainContent.classList.remove('ml-80');
    }

    hamburgerButton.addEventListener('click', () => {
        sidebar.classList.remove('-translate-x-full'); 
        hamburgerIcon.style.display = 'none';
        closeButton.style.display = 'block'; 
        mainContent.classList.add('ml-80');
    });

    closeButton.addEventListener('click', () => {
        closeSidebar();
    });

    sidebarItems.forEach(item => {
        item.addEventListener('click', closeSidebar);
    });
});