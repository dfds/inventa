using AutoMapper;
using DFDSServiceAPI.Dtos;
using Microsoft.AspNetCore.Mvc;
using Service;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace DFDSServiceAPI.Controllers
{
    [Produces("application/json")]
    [Route("api/[controller]")]
    [ApiController]
    public class ServiceProxyController : Controller
    {
        private readonly IServiceProxyService _proxyService;
        private readonly IMapper _mapper;

        public ServiceProxyController(IServiceProxyService proxyService, IMapper mapper)
        {
            _proxyService = proxyService;
            _mapper = mapper;
        }

        [HttpGet]
        [ProducesResponseType(200)]
        [ProducesResponseType(404)]
        public async Task<IActionResult> GetAll()
        {
            var results = _mapper.Map<List<ServiceProxyResultDto>>(await _proxyService.GetResults());
            return Ok(results);
        }

        [HttpGet("{k8sNamespace}")]
        public async Task<IActionResult> Get(string k8sNamespace)
        {
            var results = _mapper.Map<List<ServiceProxyResultDto>>(await _proxyService.GetResults());
            var payload = results.Select(result =>
            {
                return result.GetByNamespace(k8sNamespace);
            });

            return Ok(payload);
        }

        [HttpGet("stats")]
        [ProducesResponseType(200)]
        [ProducesResponseType(404)]
        public async Task<IActionResult> GetAllStats()
        {
            var results = _mapper.Map<List<ServiceProxyResultDto>>(await _proxyService.GetResults());
            var payload = ServiceProxyStatResultDto.FromServiceProxyResult(results);

            return Ok(payload);
        }
    }
}
