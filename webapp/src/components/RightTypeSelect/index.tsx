import * as React from 'react'
import { Select } from 'antd';
import UiStore from '../../stores/UiStore';
const Option = Select.Option;

const uiStore = UiStore.getInstance( );
export default class RightTypeSelect extends React.Component {
  render() {
    return <Select>
      <Option value={0}>{uiStore.T('RIGHT_TYPE_STANDARD')}</Option>
      <Option value={1}>{uiStore.T('RIGHT_TYPE_ADMIN')}</Option>
    </Select>
  }
}