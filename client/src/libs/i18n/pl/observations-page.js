export default {
  table: {
    title: 'Obserwacje',
    head: {
      select: 'Wybierz',
      name: 'Nazwa',
      started: 'Uruchomiona',
      lastCheckAt: 'Ostatnie sprawdzenie',
      action: 'Akcja'
    },
    started_true: 'Tak',
    started_false: 'Nie'
  },
  selected_0: `Wybrano {{count}} obserwację.`,
  selected_1: `Wybrano {{count}} obserwację.`,
  selected_2: `Wybrano {{count}} obserwacji.`,
  deleted_0: `Usunięto {{count}} obserwację.`,
  deleted_1: `Usunięto {{count}} obserwację.`,
  deleted_2: `Usunięto {{count}} obserwacji.`,
  delete: 'Usuń',
  cancel: 'Anuluj',
  created: 'Pomyślnie utworzono obserwację.',
  updated: 'Pomyślnie zaktualizowano obserwację.',
  enhancedToolbar: {
    addObservation: 'Dodaj obserwację'
  },
  observationFormDialog: {
    title: {
      edit: 'Edytowanie',
      add: 'Tworzenie'
    },
    error: {
      required: 'Wymagane pole.',
      invalidDate: 'Niepoprawna data.'
    },
    inputLabel: {
      name: 'Nazwa',
      url: 'URL',
      last_check_at: 'Od kiedy ma zacząć sprawdzać',
      started: 'Włączona',
      keyword_type: '{{index}}. Typ',
      keyword_for: '{{index}}. Dla',
      keyword_value: '{{index}}. Słowo',
      keyword_group: '{{index}}. Grupa'
    },
    keyword_title: 'Słowo #{{index}}',
    required: 'Wymagane',
    one_of: 'Jedno z',
    excluded: 'Wykluczone',
    for_title: 'Tytuł',
    description: 'Opis',
    select: 'Wybierz',
    button: {
      add: 'Dodaj',
      save: 'Zapisz',
      addKeyword: 'Dodaj słowo kluczowe'
    }
  }
};
