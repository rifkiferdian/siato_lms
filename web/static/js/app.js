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

  Alpine.data('participantPicker', function () {
    return { participantYear: '' };
  });
});

document.addEventListener('submit', function (event) {
  const form = event.target;
  if (!form.matches('[data-confirm]')) return;
  event.preventDefault();
  Swal.fire({title: form.dataset.confirm || 'Apakah Anda yakin?', icon: 'warning', showCancelButton: true, confirmButtonText: 'Ya, lanjutkan', cancelButtonText: 'Batal'}).then(result => { if (result.isConfirmed) form.submit(); });
});

document.addEventListener('DOMContentLoaded', function () {
  document.querySelectorAll('[data-participant-picker]').forEach(function (form) {
    const year = form.querySelector('[data-participant-year]');
    const search = form.querySelector('[data-participant-search]');
    const results = form.querySelector('[data-participant-results]');
    const cards = Array.from(form.querySelectorAll('[data-participant-card]'));
    const visibleCount = form.querySelector('[data-visible-count]');
    const selectedCount = form.querySelector('[data-selected-count]');
    const saveCount = form.querySelector('[data-save-count]');
    const empty = form.querySelector('[data-empty-results]');

    function updateSelection() {
      const count = cards.filter(card => card.querySelector('input').checked).length;
      selectedCount.textContent = count + ' dipilih';
      saveCount.textContent = count;
    }

    function filterStudents() {
      const selectedYear = year.value;
      const query = search.value.trim().toLowerCase();
      results.hidden = !selectedYear;
      search.disabled = !selectedYear;
      let count = 0;
      cards.forEach(function (card) {
        const matches = selectedYear !== '' && card.dataset.year === selectedYear && card.dataset.search.toLowerCase().includes(query);
        card.hidden = !matches;
        if (matches) count++;
      });
      visibleCount.textContent = count;
      empty.hidden = count !== 0;
    }

    year.addEventListener('change', function () {
      search.value = '';
      filterStudents();
    });
    search.addEventListener('input', filterStudents);
    form.querySelector('[data-select-visible]').addEventListener('click', function () {
      cards.filter(card => !card.hidden).forEach(card => { card.querySelector('input').checked = true; });
      updateSelection();
    });
    form.querySelector('[data-clear-selection]').addEventListener('click', function () {
      cards.forEach(card => { card.querySelector('input').checked = false; });
      updateSelection();
    });
    cards.forEach(card => card.querySelector('input').addEventListener('change', updateSelection));
    form.addEventListener('submit', function (event) {
      if (!cards.some(card => card.querySelector('input').checked)) {
        event.preventDefault();
        Swal.fire({icon: 'info', title: 'Pilih mahasiswa', text: 'Centang minimal satu mahasiswa sebelum menyimpan.'});
      }
    });
    filterStudents();
    updateSelection();
  });
});
