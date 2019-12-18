import React from 'react';
import PropTypes from 'prop-types';
import { List, ListItem } from '../List';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Radio from '@material-ui/core/Radio';
import { useTranslation } from 'react-i18next';
import Dialog from '../Dialog';
import { languages, fallbackLanguage } from '../../i18n';

export default function LanguageDialog({ onClose }) {
  const { t, i18n } = useTranslation();

  const codes = Object.keys(languages);
  const hasTranslations = codes.includes(i18n.language);

  return (
    <Dialog title={t('choose-language')} onClose={onClose}>
      <List>
        {codes.map(code => (
          <ListItem key={code} onClick={() => i18n.changeLanguage(code)}>
            <ListItemIcon>
              <Radio
                color="primary"
                onChange={() => i18n.changeLanguage(code)}
                checked={code === i18n.language || (!hasTranslations && code === fallbackLanguage)}
              />
            </ListItemIcon>
            <ListItemText primary={languages[code].displayName} secondary={code} />
          </ListItem>
        ))}
      </List>
    </Dialog>
  );
}

LanguageDialog.propTypes = {
  onClose: PropTypes.func.isRequired,
};
