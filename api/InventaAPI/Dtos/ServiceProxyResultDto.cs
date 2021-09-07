using k8s.Models;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Service.Classes;

namespace DFDSServiceAPI.Dtos
{
    public class ServiceProxyResultDto
    {
        public string proxyName { get; set; }
        public List<Extensionsv1beta1Ingress> ingresses { get; set; }

        public List<V1APIService> services { get; set; }

        public ServiceProxyResultDto()
        {
            ingresses = new List<Extensionsv1beta1Ingress>();
            services = new List<V1APIService>();
        }

        public ServiceProxyResultDto GetByNamespace(string k8sNamespace)
        {
            var item = new ServiceProxyResultDto();

            foreach (var ingress in ingresses)
            {
                if (ingress.Metadata.Namespace().Equals(k8sNamespace))
                {
                    item.ingresses.Add(ingress);
                }
            }
                
            // Loop through every Service object and group them by namespace
            foreach (var service in services)
            {
                if (service.Metadata.Namespace().Equals(k8sNamespace))
                {
                    item.services.Add(service);
                }
            }
            
            return item;
        }
    }

    public class ServiceProxyStatResultDto
    {
        public Dictionary<string, ServiceProxyStatItemResultDto> proxyItems { get; set; }
        public Dictionary<string, ServiceProxyStatItemResultDto> namespaceItems { get; set; }
        public ServiceProxyStatItemResultDto combinedStats { get; set; }

        public static ServiceProxyStatResultDto FromServiceProxyResult(List<ServiceProxyResultDto> results)
        {
            var payload = new ServiceProxyStatResultDto();
            payload.proxyItems = new Dictionary<string, ServiceProxyStatItemResultDto>();
            payload.namespaceItems = new Dictionary<string, ServiceProxyStatItemResultDto>();
            payload.combinedStats = new ServiceProxyStatItemResultDto();
            
            // Loop through every Inventa-operator instance
            foreach (var result in results)
            {
                var entry = new ServiceProxyStatItemResultDto
                {
                    ingressCount = result.ingresses.Count,
                    serviceCount = result.services.Count
                };
                payload.proxyItems.Add(result.proxyName, entry);

                payload.combinedStats.ingressCount += entry.ingressCount;
                payload.combinedStats.serviceCount += entry.serviceCount;
                
                // Loop through every Ingress object and group them by namespace
                foreach (var ingress in result.ingresses)
                {
                    if (!payload.namespaceItems.ContainsKey(ingress.Metadata.Namespace()))
                    {
                        payload.namespaceItems.Add(ingress.Metadata.Namespace(), new ServiceProxyStatItemResultDto());
                    }

                    payload.namespaceItems[ingress.Metadata.Namespace()].ingressCount += 1;
                }
                
                // Loop through every Service object and group them by namespace
                foreach (var service in result.services)
                {
                    if (!payload.namespaceItems.ContainsKey(service.Metadata.Namespace()))
                    {
                        payload.namespaceItems.Add(service.Metadata.Namespace(), new ServiceProxyStatItemResultDto());
                    }

                    payload.namespaceItems[service.Metadata.Namespace()].serviceCount += 1;
                }
            }

            return payload;
        }
    }

    public class ServiceProxyStatItemResultDto
    {
        public int ingressCount { get; set; }
        public int serviceCount { get; set; }
    }
}
