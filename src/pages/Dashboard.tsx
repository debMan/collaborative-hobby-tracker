import Header from '../components/layout/Header';
import Sidebar from '../components/layout/Sidebar';
import ItemList from '../components/items/ItemList';
import DetailPanel from '../components/items/DetailPanel';
import ImportModal from '../components/modals/ImportModal';
import { useStore } from '../store';

export default function Dashboard() {
  const { isDetailPanelOpen } = useStore();

  return (
    <div className="h-screen flex flex-col">
      <Header />
      <div className="flex-1 flex overflow-hidden">
        <Sidebar />
        <ItemList />
        {isDetailPanelOpen && <DetailPanel />}
      </div>
      <ImportModal />
    </div>
  );
}
