export default {
  appSettingsForm: {
    title: 'Ogólne ustawienia',
    subtitle: {
      '1': 'Wybór języka',
      '2': 'Tryb debugowania'
    },
    inputLabel: {
      lang: 'Język',
      debug: 'Włączone'
    }
  },
  collyForm: {
    title: 'Colly',
    inputLabel: {
      limit: 'Liczba requestów per IP',
      delay: 'Maksymalny czas przerwy między requestami'
    },
    error: {
      required: 'To pole jest wymagane.'
    }
  },
  discordNotificationsForm: {
    title: 'Powiadomienia Discord',
    error: {
      invalidToken: 'Nieprawidłowy token.',
      invalidChannelID: 'Nieprawidłowe ID kanału.'
    },
    inputLabel: {
      token: 'Token',
      channelID: 'ID kanału',
      enabled: 'Włączone'
    }
  },
  proxyForm: {
    title: 'Proxy',
    error: {
      invalidProxyAddress: 'Niepoprawny adres serwera proxy.'
    },
    button: {
      addProxy: 'Dodaj następne proxy'
    }
  },
  success: 'Pomyślnie zapisano zmiany.',
  error: 'Nie udało się zapisać zmian. Proszę spróbować później!',
  save: 'Zapisz',
  reset: 'Zresetuj'
};
