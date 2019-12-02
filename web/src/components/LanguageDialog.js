import React from 'react';
import PropTypes from 'prop-types';
import { List, ListItem } from './List';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Radio from '@material-ui/core/Radio';
import { useTranslation } from 'react-i18next';

import ConfigurationDialog from './ConfigurationDialog';
import { languages } from '../i18n';

export default function LanguageDialog(props) {
  const { open, onClose } = props;

  const { t, i18n } = useTranslation();

  return (
    <ConfigurationDialog
      title={t('choose-language')}
      open={open}
      onClose={onClose}
    >
      <List>
        {Object.keys(languages).map(code => (
          <ListItem key={code} onClick={() => i18n.changeLanguage(code)}>
            <ListItemIcon>
              <Radio
                color="primary"
                onChange={() => i18n.changeLanguage(code)}
                checked={i18n.language === code}
              />
            </ListItemIcon>
            <ListItemText primary={languages[code].displayName} secondary={code} />
          </ListItem>
        ))}
      </List>
    </ConfigurationDialog>
  );
}

LanguageDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
};
