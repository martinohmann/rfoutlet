import React from 'react';
import PropTypes from 'prop-types';
import Frame from './Frame';
import { Route, useRouteMatch, useHistory } from 'react-router-dom';
import { List, NoItemsListItem } from './List';
import { useTranslation } from 'react-i18next';
import GroupList from './outlets/GroupList';
import OutletScheduleDialog from './schedule/OutletScheduleDialog';
import IntervalSettingsDialog from './schedule/IntervalSettingsDialog';
import SettingsDialog from './settings/SettingsDialog';
import LanguageDialog from './settings/LanguageDialog';

export default function Routes({ groups, ready }) {
  const { t } = useTranslation();

  return (
    <Frame>
      {ready ? (
        <>
          <Route path="/">
            <GroupList groups={groups} />
          </Route>
          <Route path="/settings">
            <SettingsRoutes />
          </Route>
          <Route path="/schedule/:outletId">
            <ScheduleRoutes />
          </Route>
        </>
      ) : (
        <List>
          <NoItemsListItem
            primary={t('loading-primary')}
            secondary={t('loading-secondary')}
          />
        </List>
      )}
    </Frame>
  );
}

Routes.propTypes = {
  groups: PropTypes.array,
  ready: PropTypes.bool.isRequired,
};

const SettingsRoutes = () => {
  const match = useRouteMatch();
  const history = useHistory();

  return (
    <>
      <Route path={match.path}>
        <SettingsDialog onClose={() => history.push('/')}/>
      </Route>
      <Route path={`${match.path}/language`}>
        <LanguageDialog onClose={() => history.push(match.url)}/>
      </Route>
    </>
  );
};

const ScheduleRoutes = () => {
  const match = useRouteMatch();
  const history = useHistory();

  return (
    <>
      <Route path={match.path}>
        <OutletScheduleDialog onClose={() => history.push('/')}/>
      </Route>
      <Route path={`${match.path}/interval/:intervalId`}>
        <IntervalSettingsDialog onClose={() => history.push(match.url)}/>
      </Route>
    </>
  );
};
