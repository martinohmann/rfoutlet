import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

i18n
  // .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      en: {
        translation: {
          'settings': 'Settings',
          'schedule': 'Schedule',
          'select-weekdays': 'Select Weekdays',
          'done': 'Done',
          'add-interval': 'Add Interval',
          'edit-interval': 'Edit Interval',
          'save': 'Save',
          'edit': 'Edit',
          'delete': 'Delete',
          'weekdays': 'Weekdays',
          'start-time': 'Start Time',
          'end-time': 'End Time',
          'no-groups-primary': 'No groups configured',
          'no-groups-secondary': 'Check your rfoutlet config',
          'no-outlets-primary': 'No outlets configured for this group',
          'no-outlets-secondary': 'Check your rfoutlet config',
          'no-intervals-primary': 'No intervals configured yet',
          'no-intervals-secondary': "Tap '+' to create one",
          'loading-primary': 'Please wait.',
          'loading-secondary': 'Loading status...',
          'sunday': 'Sunday',
          'monday': 'Monday',
          'tuesday': 'Tuesday',
          'wednesday': 'Wednesday',
          'thursday': 'Thursday',
          'friday': 'Friday',
          'saturday': 'Saturday',
          'sun': 'Sun',
          'mon': 'Mon',
          'tue': 'Tue',
          'wed': 'Wed',
          'thu': 'Thu',
          'fri': 'Fri',
          'sat': 'Sat',
          'intervals-scheduled': '{{count}} interval scheduled',
          'intervals-scheduled_plural': '{{count}} intervals scheduled',
          'unset': 'unset',
        }
      },
      de: {
        translation: {
          'settings': 'Einstellungen',
          'schedule': 'Zeitschaltuhr',
          'select-weekdays': 'Wochentage wählen',
          'done': 'Fertig',
          'add-interval': 'Interval hinzufügen',
          'edit-interval': 'Interval bearbeiten',
          'save': 'Speichern',
          'edit': 'Bearbeiten',
          'delete': 'Löschen',
          'weekdays': 'Wochentage',
          'start-time': 'Startzeit',
          'end-time': 'Endzeit',
          'no-groups-primary': 'Kein Gruppen konfiguriert.',
          'no-groups-secondary': 'Überprüfe deine rfoutlet Konfiguration.',
          'no-outlets-primary': 'Keine Steckdosen für diese Gruppe konfiguriert.',
          'no-outlets-secondary': 'Überprüfe deine rfoutlet Konfiguration.',
          'no-intervals-primary': 'No intervals configured yet.',
          'no-intervals-secondary': "Drücke auf '+' zum Erstellen.",
          'loading-primary': 'Bitte warten.',
          'loading-secondary': 'Lade Status...',
          'sunday': 'Sonntag',
          'monday': 'Montag',
          'tuesday': 'Dienstag',
          'wednesday': 'Mittwoch',
          'thursday': 'Donnerstag',
          'friday': 'Freitag',
          'saturday': 'Samstag',
          'sun': 'So',
          'mon': 'Mo',
          'tue': 'Di',
          'wed': 'Mi',
          'thu': 'Do',
          'fri': 'Fr',
          'sat': 'Sa',
          'intervals-scheduled': '{{count}} Interval konfiguriert',
          'intervals-scheduled_plural': '{{count}} Intervalle konfiguriert',
          'unset': 'nicht gesetzt',
        }
      }
    },
    lng: 'de',
    fallbackLng: 'en',
    debug: true,
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
