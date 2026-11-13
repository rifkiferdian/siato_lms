document.addEventListener('alpine:init', function () {
  Alpine.data('dashboardShell', function () {
    return {
      sidebar: false,
      profile: false,
      openSidebar() { this.sidebar = true; },
      closeSidebar() { this.sidebar = false; }
    };
  });

  Alpine.data('loginForm', function () {
    return { showPassword: false };
  });

  Alpine.data('schedulePage', function () {
    return { scheduleModal: false };
  });

  Alpine.data('classPage', function () {
    return { classModal: false };
  });

  Alpine.data('peoplePage', function () {
    return { personModal: false };
  });

  Alpine.data('userManagementPage', function () {
    return { createUserModal: false };
  });
});

document.addEventListener('submit', function (event) {
  const form = event.target;
  if (!form.matches('[data-confirm]')) return;
  event.preventDefault();
  Swal.fire({title: form.dataset.confirm || 'Apakah Anda yakin?', icon: 'warning', showCancelButton: true, confirmButtonText: 'Ya, lanjutkan', cancelButtonText: 'Batal'}).then(result => { if (result.isConfirmed) form.submit(); });
});
