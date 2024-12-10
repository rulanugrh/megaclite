document.addEventListener('DOMContentLoaded', function() {
    // Get references to elements
    const sidebar = document.getElementById('sidebar');
    const hamburgerButton = document.getElementById('hamburgerButton');
    const hamburgerIcon = document.getElementById('hamburgerIcon');
    const closeButton = document.getElementById('closeButton');
    const mainContent = document.getElementById('mainContent');

    // Get all sidebar menu items
    const sidebarItems = document.querySelectorAll('#mail a');
    
    // Function to close sidebar
    function closeSidebar() {
        sidebar.classList.add('-translate-x-full'); // Hide sidebar
        hamburgerIcon.style.display = 'block'; // Show hamburger icon again
        if (window.innerWidth >= 1024) { // Large screen
            closeButton.style.display = 'block'; // Always display close button on large screens
        } else { // Small screen
            closeButton.style.display = 'none'; // Hide close button on small screens
        }
        mainContent.classList.remove('ml-80'); // Restore main content layout
    }

    // Toggle the sidebar visibility when the hamburger button is clicked
    hamburgerButton.addEventListener('click', () => {
        sidebar.classList.remove('-translate-x-full'); // Show sidebar
        hamburgerIcon.style.display = 'none'; // Hide hamburger icon
        closeButton.style.display = 'block'; // Show close button (X)
        mainContent.classList.add('ml-80'); // Shift the main content when sidebar is open
    });

    // Close the sidebar when the close button (X) is clicked
    closeButton.addEventListener('click', () => {
        closeSidebar(); // Use the closeSidebar function
    });

    // Close the sidebar when a sidebar menu item is clicked
    sidebarItems.forEach(item => {
        item.addEventListener('click', closeSidebar); // Close the sidebar when any item is clicked
    });
});