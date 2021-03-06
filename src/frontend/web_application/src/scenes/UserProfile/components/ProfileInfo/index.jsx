import { createSelector } from 'reselect';
import { compose } from 'redux';
import { connect } from 'react-redux';
import { withTranslator } from '@gandi/react-translate';
import Presenter from './presenter';

const i18nSelector = state => state.i18n;
const mapStateToProps = createSelector(
  [i18nSelector],
  i18n => ({
    locale: i18n.locale,
  })
);

export default compose(
  connect(mapStateToProps),
  withTranslator()
)(Presenter);
