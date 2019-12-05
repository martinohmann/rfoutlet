import React from 'react';
import PropTypes from 'prop-types';
import { List, ListItem } from './List';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import LanguageIcon from '@material-ui/icons/Language';
import { useTranslation } from 'react-i18next';
import { Route, useHistory, useRouteMatch } from 'react-router';

import ConfigurationDialog from './ConfigurationDialog';
import LanguageDialog from './LanguageDialog';
import GitHubIcon from './GitHubIcon';
import config from '../config';

export default function SettingsDialog(props) {
  const history = useHistory();
  const { path, url } = useRouteMatch();

  const handleDialogOpen = (open) => () => history.push(open ? `${url}/language` : url);

  const { onClose } = props;

  const { t, i18n } = useTranslation();

  return (
    <ConfigurationDialog
      title={t('settings')}
      onClose={onClose}
    >
      <List>
        <ListItem onClick={handleDialogOpen(true)}>
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
      <Route path={`${path}/language`}>
        <LanguageDialog onClose={handleDialogOpen(false)} />
      </Route>
    </ConfigurationDialog>
  );
}

SettingsDialog.propTypes = {
  onClose: PropTypes.func.isRequired,
};
