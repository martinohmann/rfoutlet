import React from 'react';
import PropTypes from 'prop-types';
import { List, NoItemsListItem } from './List';
import GroupListItem from './GroupListItem';
import { useTranslation } from 'react-i18next';

export default function GroupList(props) {
  const { groups } = props;
  const { t } = useTranslation();

  return (
    <List>
      {groups.map(group =>
        <GroupListItem key={group.id} {...group} />
      )}
      {groups.length === 0 ? (
        <NoItemsListItem
          primary={t('no-groups-primary')}
          secondary={t('no-groups-secondary')}
        />
      ) : ''}
    </List>
  );
}

GroupList.propTypes = {
  groups: PropTypes.array.isRequired,
};
