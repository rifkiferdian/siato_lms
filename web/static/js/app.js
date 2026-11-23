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
  if (window.location.pathname === '/admin/courses') {
    const courseTable = document.querySelector('.people-table');
    if (courseTable) {
      const page = Math.max(1, Number(new URLSearchParams(window.location.search).get('page')) || 1);
      const heading = document.createElement('th');
      heading.textContent = 'No.';
      heading.className = 'w-16 text-center';
      courseTable.querySelector('thead tr').prepend(heading);
      courseTable.querySelectorAll('tbody tr').forEach(function (row, index) {
        const cell = document.createElement('td');
        cell.textContent = String((page - 1) * 20 + index + 1);
        cell.className = 'text-center font-semibold text-slate-400';
        row.prepend(cell);
      });
    }

    const courseModalForm = document.querySelector('#course-modal form');
    if (courseModalForm && !courseModalForm.querySelector('[name="theory_credits"]')) {
      const typeRow = courseModalForm.querySelector('[name="course_type"]').closest('.grid');
      const componentFields = document.createElement('div');
      componentFields.className = 'course-component-fields';
      componentFields.innerHTML = '<div class="mb-3 flex items-center justify-between"><div><p class="text-xs font-semibold text-slate-700">Komponen SKS</p><p class="mt-1 text-[11px] text-slate-400">Jumlah komponen harus sama dengan total SKS.</p></div><span class="rounded-lg bg-blue-50 px-2.5 py-1 text-[10px] font-bold text-blue-600">T + P + PR + L</span></div><div class="grid grid-cols-2 gap-3 sm:grid-cols-4"><label><span class="form-label">Teori (T)</span><input required type="number" name="theory_credits" min="0" max="12" value="0" class="form-input"></label><label><span class="form-label">Praktikum (P)</span><input required type="number" name="practicum_credits" min="0" max="12" value="0" class="form-input"></label><label><span class="form-label">Praktik (PR)</span><input required type="number" name="practice_credits" min="0" max="12" value="0" class="form-input"></label><label><span class="form-label">Lapangan (L)</span><input required type="number" name="field_credits" min="0" max="12" value="0" class="form-input"></label></div><p data-component-total class="mt-2 text-right text-[11px] font-semibold text-slate-500">Total komponen: 0 SKS</p>';
      typeRow.parentNode.insertBefore(componentFields, typeRow);

      function updateComponentTotal() {
        const total = ['theory_credits','practicum_credits','practice_credits','field_credits'].reduce(function (sum, name) { return sum + (Number(courseModalForm.elements[name].value) || 0); }, 0);
        const expected = Number(courseModalForm.elements.credits.value) || 0;
        const output = componentFields.querySelector('[data-component-total]');
        output.textContent = 'Total komponen: ' + total + ' SKS';
        output.className = 'mt-2 text-right text-[11px] font-semibold ' + (total === expected ? 'text-emerald-600' : 'text-amber-600');
      }
      courseModalForm.addEventListener('input', updateComponentTotal);
      courseModalForm.addEventListener('change', updateComponentTotal);
    }

    document.querySelectorAll('[data-academic-edit][data-modal="course-modal"]').forEach(function (button) {
      if (button.dataset.action === '/admin/courses') {
        Object.assign(button.dataset, {credits:'3',recommendedSemester:'1',courseType:'required',isActive:'1',theoryCredits:'3',practicumCredits:'0',practiceCredits:'0',fieldCredits:'0'});
        return;
      }
      const match = button.dataset.action.match(/\/admin\/courses\/(\d+)$/);
      const stored = match && document.querySelector('[data-course-component-store] [data-course-id="' + match[1] + '"]');
      if (stored) Object.assign(button.dataset, {theoryCredits:stored.dataset.theoryCredits,practicumCredits:stored.dataset.practicumCredits,practiceCredits:stored.dataset.practiceCredits,fieldCredits:stored.dataset.fieldCredits});
    });
  }

  const academicModals = document.querySelectorAll('[data-academic-modal]');

  function closeAcademicModal(modal) {
    if (!modal) return;
    modal.hidden = true;
    document.body.classList.remove('overflow-hidden');
  }

  document.querySelectorAll('[data-academic-edit]').forEach(function (button) {
    button.addEventListener('click', function () {
      const modal = document.getElementById(button.dataset.modal);
      if (!modal) return;
      const form = modal.querySelector('form');
      form.reset();
      form.action = button.dataset.action;
      form.querySelectorAll('[name]').forEach(function (field) {
        if (field.name === '_csrf') return;
        const key = field.name.replace(/_([a-z])/g, function (_, letter) { return letter.toUpperCase(); });
        if (Object.prototype.hasOwnProperty.call(button.dataset, key)) field.value = button.dataset[key];
      });
      form.dispatchEvent(new Event('input', {bubbles: true}));
      modal.hidden = false;
      document.body.classList.add('overflow-hidden');
      const firstField = form.querySelector('input:not([type="hidden"]), select');
      if (firstField) firstField.focus();
    });
  });

  academicModals.forEach(function (modal) {
    modal.querySelectorAll('[data-modal-close]').forEach(function (button) {
      button.addEventListener('click', function () { closeAcademicModal(modal); });
    });
  });

  document.addEventListener('keydown', function (event) {
    if (event.key !== 'Escape') return;
    academicModals.forEach(function (modal) { if (!modal.hidden) closeAcademicModal(modal); });
  });

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
