using k8s.Models;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace DFDSServiceUi
{
    public class ServiceProxyResult
    {
        public string proxyName { get; set; }
        public List<Extensionsv1beta1Ingress> ingresses { get; set; }

        public List<V1APIService> services { get; set; }
    }
}
