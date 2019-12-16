import React from 'react';
import PropTypes from 'prop-types';
import { List, ListItem } from '../List';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import LanguageIcon from '@material-ui/icons/Language';
import { useTranslation } from 'react-i18next';
import { useHistory, useRouteMatch } from 'react-router';
import Dialog from '../Dialog';
import SvgIcon from '@material-ui/core/SvgIcon';
import config from '../../config';

export default function SettingsIndex({ onClose }) {
  const history = useHistory();
  const { url } = useRouteMatch();
  const { t, i18n } = useTranslation();

  return (
    <Dialog title={t('settings')} onClose={onClose}>
      <List>
        <ListItem onClick={() => history.push(`${url}/language`)}>
          <ListItemIcon>
            <LanguageIcon />
          </ListItemIcon>
          <ListItemText primary={t('language-settings')} secondary={i18n.language !== 'en' ? 'Language Settings' : ''} />
        </ListItem>
        <ListItem onClick={() => window.open(config.project.url, '_blank')}>
          <ListItemIcon>
            <GitHubIcon />
          </ListItemIcon>
          <ListItemText primary={t('project-on-github')} secondary={config.project.url.replace('https://', '')} />
        </ListItem>
      </List>
    </Dialog>
  );
}

SettingsIndex.propTypes = {
  onClose: PropTypes.func.isRequired,
};

const GitHubIcon = () => {
  return (
    <SvgIcon>
      <svg focusable="false" viewBox="0 0 24 24" aria-hidden="true" role="presentation">
        <path d="M12 .3a12 12 0 0 0-3.8 23.4c.6.1.8-.3.8-.6v-2c-3.3.7-4-1.6-4-1.6-.6-1.4-1.4-1.8-1.4-1.8-1-.7.1-.7.1-.7 1.2 0 1.9 1.2 1.9 1.2 1 1.8 2.8 1.3 3.5 1 0-.8.4-1.3.7-1.6-2.7-.3-5.5-1.3-5.5-6 0-1.2.5-2.3 1.3-3.1-.2-.4-.6-1.6 0-3.2 0 0 1-.3 3.4 1.2a11.5 11.5 0 0 1 6 0c2.3-1.5 3.3-1.2 3.3-1.2.6 1.6.2 2.8 0 3.2.9.8 1.3 1.9 1.3 3.2 0 4.6-2.8 5.6-5.5 5.9.5.4.9 1 .9 2.2v3.3c0 .3.1.7.8.6A12 12 0 0 0 12 .3"></path>
      </svg>
    </SvgIcon>
  );
};
