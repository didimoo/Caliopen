import { createSelector } from 'reselect';
import { connect } from 'react-redux';
import { settingsSelector } from '../store/selectors/settings';

const mapStateToProps = createSelector(
  [settingsSelector],
  settings => ({
    settings,
  })
);

const withSettings = () => Component => connect(mapStateToProps)(Component);

export { withSettings };
