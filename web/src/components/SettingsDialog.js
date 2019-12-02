import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { List, ListItem } from './List';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import LanguageIcon from '@material-ui/icons/Language';
import { useTranslation } from 'react-i18next';

import ConfigurationDialog from './ConfigurationDialog';
import LanguageDialog from './LanguageDialog';
import GitHubIcon from './GitHubIcon';
import config from '../config';

export default function SettingsDialog(props) {
  const [languageDialogOpen, setLanguageDialogOpen] = useState(false);

  const handleDialogOpen = (open) => () => setLanguageDialogOpen(open);

  const { open, onClose } = props;

  const { t, i18n } = useTranslation();

  return (
    <ConfigurationDialog
      title={t('settings')}
      open={open}
      onClose={onClose}
    >
      <List>
        <ListItem onClick={handleDialogOpen(true)}>
          <ListItemIcon>
            <LanguageIcon />
          </ListItemIcon>
          <ListItemText primary={t('language-settings')} secondary={i18n.language !== 'en' ? 'Language Settings' : ''} />
        </ListItem>
        <ListItem onClick={() => window.location.href = config.project.url}>
          <ListItemIcon>
            <GitHubIcon />
          </ListItemIcon>
          <ListItemText primary={t('project-on-github')} secondary={config.project.url.replace('https://', '')} />
        </ListItem>
      </List>
      <LanguageDialog
        open={languageDialogOpen}
        onClose={handleDialogOpen(false)}
      />
    </ConfigurationDialog>
  );
}

SettingsDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
};
