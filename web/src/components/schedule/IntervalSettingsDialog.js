import React, { useReducer } from 'react';
import PropTypes from 'prop-types';
import { useTranslation } from 'react-i18next';
import CheckIcon from '@material-ui/icons/Check';
import { Route, Switch, useHistory, useRouteMatch } from 'react-router';
import { useCurrentInterval, useCurrentOutlet } from '../../hooks';
import Dialog from '../Dialog';
import IntervalSettingsList from './IntervalSettingsList';
import IntervalTimePicker from './IntervalTimePicker';
import WeekdaysDialog from './WeekdaysDialog';
import dispatcher from '../../dispatcher';

const emptyInterval = {
  id: null,
  enabled: false,
  from: null,
  to: null,
  weekdays: [],
};

const reduceState = (state, changes) => ({ ...state, ...changes });

export default function IntervalSettingsDialog({ onClose }) {
  const history = useHistory();
  const { path, url } = useRouteMatch();
  const outlet = useCurrentOutlet();
  const interval = useCurrentInterval() || emptyInterval;

  const [state, setState] = useReducer(reduceState, interval);

  const handleOpen = (name, open) => () => history.push(open ? `${url}/${name}` : url);

  const handleChange = (name) => (value) => setState({ [name]: value });

  const handleDone = () => {
    const newInterval = {
      ...interval,
      weekdays: state.weekdays,
      from: state.from,
      to: state.to
    };

    const action = newInterval.id ? 'update' : 'create';

    dispatcher.dispatchIntervalMessage(outlet.id, action, newInterval);

    onClose();
  }

  const isComplete = () => state.from && state.to && state.weekdays.length > 0;

  const { t } = useTranslation();

  return (
    <Dialog
      title={state.id ? t('edit-interval') : t('add-interval')}
      onClose={onClose}
      onDone={handleDone}
      doneButtonDisabled={!isComplete()}
      doneButtonText={<CheckIcon />}
    >
      <IntervalSettingsList
        weekdays={state.weekdays}
        fromDayTime={state.from}
        toDayTime={state.to}
        onWeekdaysClick={handleOpen('weekdays', true)}
        onFromDayTimeClick={handleOpen('from', true)}
        onToDayTimeClick={handleOpen('to', true)}
      />
      <Switch>
        <Route path={`${path}/from`}>
          <IntervalTimePicker
            value={state.from}
            onChange={handleChange('from')}
            onClose={handleOpen('from', false)}
          />
        </Route>
        <Route path={`${path}/to`}>
          <IntervalTimePicker
            value={state.to}
            onChange={handleChange('to')}
            onClose={handleOpen('to', false)}
          />
        </Route>
        <Route path={`${path}/weekdays`}>
          <WeekdaysDialog
            selected={state.weekdays}
            key={interval.id}
            onChange={handleChange('weekdays')}
            onClose={handleOpen('weekdays', false)}
          />
        </Route>
      </Switch>
    </Dialog>
  );
}

IntervalSettingsDialog.propTypes = {
  onClose: PropTypes.func.isRequired,
};
